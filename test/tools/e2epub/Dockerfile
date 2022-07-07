FROM golang:1.18 AS build
WORKDIR /go/src/
COPY . /go/src/
ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev

WORKDIR /go/src/test/tools/e2epub
RUN ls
RUN go mod tidy && GOOS=linux go build -o /install/e2epub main.go

FROM alpine:3.12
COPY --from=build /install/e2epub /e2epub
COPY ./test/tools/e2epub/entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]