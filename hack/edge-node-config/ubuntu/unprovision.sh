#!/bin/bash
yurtadm reset -f
rm -rf /etc/cni/net.d
ip link set cni0 down
ip link delete cni0
ip link set flannel.1 down
ip link delete flannel.1
ip link set yurthub-dummy0 down
ip link delete yurthub-dummy0
systemctl restart docker
