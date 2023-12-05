# EdgeFarm Edge Node running Ubuntu 22.04

This guide will give instructions on how to run a vanilla Ubuntu 22.04 installation as 
Edge Node in an EdgeFarm environment.

## Preparations

In order to run run the needed components with all the cool features you need to have some 
prerequisits to be installed:

- docker
- netbird
- socat an conntrack

Also any swap partitions and swap files should be disabled permanently.

You can use the `prepare.sh` shell script that comes with this package.

Simply run `prepare.sh` and let the script take care of the installation. Follow the steps after the script to 
finish the preparation step.

## Installation

Run the `install.sh` script with the needed arguments. See the `--help` options for possible arguments.

This command will download everything needed and configures the system. In the end it joins the Kubernetes Cluster.
```
$ ./install.sh --token <your-bootstrap-token> --address <your-address:port> --join --node-ip $(cat /usr/local/etc/wt0.ip)
```

## Unprovisioning a edge node from the cluster

To unprovision a edge node from the cluster, run the `unprovision.sh` script with the needed arguments. This will delete some network interfaces and deletes previously created containers.
After running this command you can re-provision the edge node with the `install.sh` script.
