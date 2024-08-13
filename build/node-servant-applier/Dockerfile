FROM golang:1.23 as build

WORKDIR /app
COPY . .
RUN go mod tidy && go mod vendor
RUN go build -ldflags "-linkmode external -extldflags -static" -o cmd/node-servant-applier/main cmd/node-servant-applier/main.go

FROM scratch
COPY --from=build /app/cmd/node-servant-applier/main /main
CMD ["/main"]