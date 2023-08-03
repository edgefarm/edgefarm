NAME = local-up
BIN_DIR ?= dist
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -gcflags "all=-N -l" -ldflags '-extldflags "-static"' -ldflags "-X github.com/edgefarm/edgefarm/cmd/local-up/cmd.Version=$(VERSION)"
GO_ARCH ?= amd64
GO_OS ?= linux

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

mod: ## go mod handling
	go mod tidy
	go mod vendor

test: ## run tests
	go test ./...

build: ## build local-up tool
	cd cmd/local-up && CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o ../../${BIN_DIR}/${NAME}-${GO_OS}-${GO_ARCH} main.go

clean: ## remove files created during build pipeline
	rm -rf ${BIN_DIR}

.PHONY: help all build test clean mod
