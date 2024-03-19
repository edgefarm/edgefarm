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
rm -fr /etc/kubernetes/
rm -fr ~/.kube/
rm -fr /var/lib/etcd
rm -rf /var/lib/cni/
rm -rf /etc/cni/net.d
ip link set cni0 down
ip link delete cni0
ip link set flannel.1 down
ip link delete flannel.1
ip link set yurthub-dummy0 down
ip link delete yurthub-dummy0
apt purge kubectl kubeadm kubelet kubernetes-cni cri-tools -y
systemctl daemon-reload

docker rm -f `docker ps -a | grep "k8s_" | awk '{print $1}'`
systemctl restart docker
iptables -F 
iptables -t nat -F
iptables -t mangle -F
iptables -X
