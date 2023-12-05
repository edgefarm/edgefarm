#!/bin/bash

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
DIST=${1}
TARGET=${SCRIPTPATH}/${DIST}
echo $DIST
TMP=$(mktemp -d)
cp -r ubuntu-22.04 ${TMP}/edgefarm-node-config
cd ${TMP}
ls
tar cfvz ${TARGET}/ubuntu-22.04-edge-node-config.tar.gz edgefarm-node-config
