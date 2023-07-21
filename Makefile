NAME = local-up
BIN_DIR ?= dist
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/edgefarm/edgefarm/cmd/local-up/cmd.version=$(VERSION)"
GO_ARCH = amd64


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

mod: ## go mod handling
	go mod tidy
	go mod vendor

test: ## run tests
	go test ./...

build: ## build local-up tool
	cd cmd/local-up && GOOS=linux GOARCH=${GO_ARCH} go build $(GO_LDFLAGS) -o ../../${BIN_DIR}/${NAME}-${GO_ARCH} main.go

clean: ## remove files created during build pipeline
	rm -rf ${BIN_DIR}

.PHONY: all build test clean mod