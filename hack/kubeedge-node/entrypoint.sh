#!/bin/bash
# generate kubeedge config
/kubeedge/generate.sh

mount --make-rshared /
source /opt/bash-utils/logger.sh

function wait_for_process () {
    local max_time_wait=30
    local process_name="$1"
    local waited_sec=0
    while ! pgrep "$process_name" >/dev/null && ((waited_sec < max_time_wait)); do
        INFO "Process $process_name is not running yet. Retrying in 1 seconds"
        INFO "Waited $waited_sec seconds of $max_time_wait seconds"
        sleep 1
        ((waited_sec=waited_sec+1))
        if ((waited_sec >= max_time_wait)); then
            return 1
        fi
    done
    return 0
}

INFO "Starting supervisor"
/usr/bin/supervisord -n >> /dev/null 2>&1 &

INFO "Waiting for processes to be running"
processes=(dockerd)

for process in "${processes[@]}"; do
    wait_for_process "$process"
    if [ $? -ne 0 ]; then
        ERROR "$process is not running after max time"
        exit 1
    else
        INFO "$process is running"
    fi
done

# Start the first process
echo Starting edgecore
mkdir -p /var/log/kubeedge
/usr/local/bin/edgecore
#--log_dir /var/log/kubeedge
status=$?
if [ $status -ne 0 ]; then
    echo "Failed to start my_second_process: $status"
    exit $status
fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds

while sleep 60; do
    ps aux |grep my_first_process |grep -q -v grep
    PROCESS_1_STATUS=$?
    #  ps aux |grep my_second_process |grep -q -v grep
    #  PROCESS_2_STATUS=$?
    # If the greps above find anything, they exit with 0 status
    # If they are not both 0, then something is wrong
    #  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
    if [ $PROCESS_1_STATUS -ne 0 ]; then
        echo "processes has already exited."
        exit 1
    fi
done
