all: container
ARCH?=amd64
CONTAINER_ARCHS?=linux/amd64,linux/arm64
VERSION?=$(shell git describe --match=NeVeRmAtCh --always --abbrev=8 --dirty)
REGISTRY?=ghcr.io/edgefarm/edgefarm
IMAGE_NAME?=example-basic-producer
GO_LDFLAGS = -tags 'netgo osusergo static_build'

proto:
	cd proto && make

container:
	docker buildx build --build-arg VERSION=${VERSION} -f Dockerfile \
		--platform ${CONTAINER_ARCHS} --push -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} -t ${REGISTRY}/${IMAGE_NAME}:latest .

producer:
	GOOS=linux GOARCH=${ARCH} go build $(GO_LDFLAGS) -o ./bin/producer ./cmd/producer/main.go


.PHONY: all proto container producer
