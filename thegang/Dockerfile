FROM golang:1.20 AS base

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -y && \
    apt-get install -y wget build-essential pkg-config ca-certificates --no-install-recommends make

WORKDIR /go/app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

### Development with hot reload and debugger
FROM base AS dev
WORKDIR /go/app

# Hot reloading mod
RUN go install github.com/cosmtrek/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["air"]