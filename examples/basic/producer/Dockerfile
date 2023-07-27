FROM golang:1.20 AS build
WORKDIR /go/src/
COPY . /go/src/
ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev

WORKDIR /go/src
RUN go mod tidy && make producer

FROM alpine:3.12
COPY --from=build /go/src/bin/producer /producer
RUN chmod +x /producer
ENTRYPOINT ["/producer"]
