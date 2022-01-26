#!/bin/bash
docker build -t edgefarm/virtual-device:latest -f hack/kubeedge-node/Dockerfile .
