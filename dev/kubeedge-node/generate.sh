#!/usr/bin/env bash
NODE_IP=$(hostname -i)
template=/kubeedge/edgecore.yaml.TEMPLATE
values=/kubeedge/values.conf

echo CLOUDCORE_ADDRESS=${CLOUDCORE_ADDRESS} > ${values}
echo NODE_NAME=${NODE_NAME} >> ${values}
echo NODE_IP=${NODE_IP} >> ${values}

mkdir -p /etc/kubeedge/config/
. "${values}"
eval "echo \"$(cat "${template}")\"" > /etc/kubeedge/config/edgecore.yaml
