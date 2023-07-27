SRC_DIR=.
all: ${SRC_DIR}/go/sensor.go

prepare:
	mkdir -p go/ 
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

${SRC_DIR}/go/sensor.go: ${SRC_DIR}/sensor.proto prepare
	protoc -I=${SRC_DIR} ${SRC_DIR}/sensor.proto --go_out=go/

.PHONY: all prepare
