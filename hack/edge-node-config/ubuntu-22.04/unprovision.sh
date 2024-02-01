#!/bin/bash

# Try to reset via yurtadm first
yurtadm --help > /dev/null 2>&1
if [ $? -eq 0 ]; then
  yurtadm reset -f
else
  # Try to reset via kubeadm as a fallback
  kubeadm --help > /dev/null 2>&1
  if [ $? -eq 0 ]; then
    kubeadm reset -f
  fi
fi
rm -rf /etc/cni/net.d
ip link set cni0 down
ip link delete cni0
ip link set flannel.1 down
ip link delete flannel.1
ip link set yurthub-dummy0 down
ip link delete yurthub-dummy0
systemctl restart docker
