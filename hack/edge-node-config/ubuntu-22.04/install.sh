#!/bin/bash

green='\033[0;32m'
red='\033[0;31m'
nc='\033[0m'

###### PARAMETERS HANDLING
NAME=${NAME:-$(hostname)}
TOKEN=${TOKEN:-}
ARCH=${ARCH:-$(uname -m)}
ADDRESS=${ADDRESS:-}
JOIN=${JOIN:-false}
PRECHECKS_ONLY=${PRECHECKS_ONLY:-false}
NODE_IP=${NODE_IP:-}

options=$(getopt -o "" -l "prechecks-only,join,node-ip:,address:,name:,token:,arch:,help" -- "$@")

if [ $? -ne 0 ]; then
 echo "Invalid arguments."
 exit 1
fi

eval set -- "$options"
Help() {
 echo "Usage: script --address ADDRESS --token TOKEN [--name NAME] [--arch ARCH] [--node-ip IP] [--help]"
 echo
 echo "Options:"
 echo "--address ADDRESS   Set the API server address. This option is mandatory."
 echo "--token TOKEN       Set the token. This option is mandatory."
 echo "--name NAME         Set the node name. Default is the hostname. (optional)"
 echo "--arch ARCH         Set the architecture. Allowed values are arm64, amd64, and arm. Default is the current architecture. (optional)"
 echo "--node-ip IP        Set the node ip. (optional)"
 echo "--join              Run the join command rather than printing it after setting everything up. (optional)"
 echo "--prechecks-only         Run prechecks only. (optional)"
 echo "--help              Display this help message."
 echo
}

while [ $# -gt 0 ]; do
 case "$1" in
  --address) ADDRESS="$2"; shift;;
  --name) NAME="$2"; shift;;
  --token) TOKEN="$2"; shift;;
  --arch) ARCH="$2"; shift;;
  --node-ip) NODE_IP="$2"; shift;;
  --help) Help; exit;;
  --join) JOIN="true";;
  --prechecks-only) PRECHECKS_ONLY="true";;
  --) shift;;
 esac
 shift
done

if [ -z "$TOKEN" ]; then
 echo -e "${red}Token must be set.${nc}"
 exit 1
fi

if [ -z "$ADDRESS" ]; then
 echo -e "${red}Address must be set.${nc}"
 exit 1
fi

# Map uname architecture to specific values
case "$ARCH" in
 "x86_64") ARCH="amd64" ;;
 "aarch64") ARCH="arm64" ;;
 "armv7l") ARCH="arm" ;;
 *) echo -e "${red}Invalid architecture. Allowed values are arm64, amd64 and arm.${nc}"; exit 1 ;;
esac

###### PRECHECKS
PRECHECK_ERRORS=0
echo "Running prechecks..."
INTERFACE="wt0"
if ip link show "$INTERFACE" | grep -qs "state UP"; then
   echo -e "  ${red}Interface $INTERFACE is not up or does not exist. Make sure that netbird is installed, up and running.${nc}"
   PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
else
    echo -e "  ${green}Interface $INTERFACE is up.${nc}"
fi

# Try to run a Docker command
docker info > /dev/null 2>&1
# Check the exit status
if [ $? -eq 0 ]; then
   echo -e "  ${green}Docker is running${nc}"
else
   echo -e "  ${red}Docker is not running${nc}"
   PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
fi

if [ $PRECHECK_ERRORS -gt 0 ]; then
    echo -e "${red}Prechecks failed.${nc}"
    exit 1
fi
echo "Prechecks passed."

###### INSTALLATION

if $PRECHECKS_ONLY == "true" ; then
    exit 0
fi

TMP=$(mktemp -d)

mkdir -p /usr/local/bin
mkdir -p /opt/cni/bin
mkdir -p /etc/edgefarm/
mkdir -p /etc/systemd/system
mkdir -p /etc/udev/rules.d

echo "Downloading components..."
wget -q --show-progress https://github.com/openyurtio/openyurt/releases/download/v1.4.0/yurtadm-v1.4.0-linux-${ARCH}.tar.gz -P ${TMP}
tar xfz ${TMP}/yurtadm-v1.4.0-linux-${ARCH}.tar.gz  -C ${TMP} && mv ${TMP}/linux-${ARCH}/yurtadm /usr/local/bin/yurtadm && chmod +x /usr/local/bin/yurtadm

wget -q --show-progress https://github.com/edgefarm/edgefarm/releases/download/cni-0.8.0/cni-plugins-linux-${ARCH}-v0.8.0.tgz -P ${TMP}
tar xfz ${TMP}/cni-plugins-linux-${ARCH}-v0.8.0.tgz -C /opt/cni/bin --pax-option=delete=SCHILY.*,delete=LIBARCHIVE.*

cp files/kubeadm-join.conf.template ${TMP}/kubeadm-join.conf
sed -i "s/ADDRESS/$ADDRESS/g" ${TMP}/kubeadm-join.conf
sed -i "s/NODE_NAME/$NAME/g" ${TMP}/kubeadm-join.conf
sed -i "s/BOOTSTRAP_TOKEN/$TOKEN/g" ${TMP}/kubeadm-join.conf

if [ -n "${NODE_IP}" ]; then
  echo "    node-ip: ${NODE_IP}" >> ${TMP}/kubeadm-join.conf
fi

cp ${TMP}/kubeadm-join.conf /etc/edgefarm/
rm -rf ${TMP}

cp files/etc/systemd/system/edge0-device.service /etc/systemd/system/
cp files/etc/udev/rules.d/90-wt0.rules /etc/udev/rules.d/
cp files/usr/local/bin/add-edge0.sh /usr/local/bin/
cp files/usr/local/bin/add-wt0.sh /usr/local/bin/
cp files/usr/local/bin/remove-wt0.sh /usr/local/bin/

systemctl enable edge0-device
systemctl start edge0-device
udevadm control --reload-rules
udevadm trigger
/usr/local/bin/add-wt0.sh

###### JOIN CLUSTER
if $JOIN ; then
    echo -e "${green}Joining the cluster...${nc}"
    yurtadm join ${ADDRESS} --config /etc/edgefarm/kubeadm-join.conf --node-name=${NAME} --token=${TOKEN} --node-type=edge --discovery-token-unsafe-skip-ca-verification --v=9 --reuse-cni-bin --yurthub-image ghcr.io/openyurtio/openyurt/yurthub:v1.4.0 --cri-socket /var/run/dockershim.sock --yurthub-server-addr=192.168.168.1
else
    echo -e "${green}Everything is set up. Run the following command to join the cluster:${nc}"
    echo yurtadm join ${ADDRESS} --config /etc/edgefarm/kubeadm-join.conf --node-name=${NAME} --token=${TOKEN} --node-type=edge --discovery-token-unsafe-skip-ca-verification --v=9 --reuse-cni-bin --yurthub-image ghcr.io/openyurtio/openyurt/yurthub:v1.4.0 --cri-socket /var/run/dockershim.sock --yurthub-server-addr=192.168.168.1
fi
