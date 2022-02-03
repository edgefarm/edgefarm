#!/bin/sh
#
# Execute this script for provisioning edgefarm certs to a device.
#
# To be called with arguments: <device ip address>
#
# Requires generated certificates in folder <repository root>/dev/manifests/kubeedge-certs/config/
# - rootCa.pem
# - node.pem
# - node.key
#

if [ "$#" -ne 1 ] ; then
    echo "Usage: ${0} <device ip address>"
    exit 1
fi

deviceIP=${1}

scp dev/manifests/kubeedge-certs/config/rootCa.pem root@${deviceIP}:/etc/kubeedge/certs/
scp dev/manifests/kubeedge-certs/config/node.pem root@${deviceIP}:/etc/kubeedge/certs/
scp dev/manifests/kubeedge-certs/config/node.key root@${deviceIP}:/etc/kubeedge/certs/
