all: container
CONTAINER_ARCHS?=linux/amd64,linux/arm64
REGISTRY?=ghcr.io/edgefarm/edgefarm
IMAGE_NAME?=example-basic-consumer
VERSION?=$(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)

proto:
	cd proto && make

container:
	docker buildx build --build-arg VERSION=${VERSION} -f Dockerfile \
		--platform ${CONTAINER_ARCHS} --push -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} -t ${REGISTRY}/${IMAGE_NAME}:latest .


.PHONY: all proto container
