#!/bin/bash
CLOUDCORE_ADDRESS=${1}
NODE_NAME=${2}
docker run -d --rm --env CLOUDCORE_ADDRESS=${CLOUDCORE_ADDRESS} --env NODE_NAME=${NODE_NAME} --name ${NODE_NAME} \
-v $(pwd)/dev/manifests/kubeedge-certs/config/rootCa.pem:/etc/kubeedge/certs/rootCa.pem \
-v $(pwd)/dev/manifests/kubeedge-certs/config/node.pem:/etc/kubeedge/certs/node.pem \
-v $(pwd)/dev/manifests/kubeedge-certs/config/node.key:/etc/kubeedge/certs/node.key \
--privileged edgefarm/virtual-device:latest