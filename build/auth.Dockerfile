FROM golang:1.24.0-alpine AS builder

COPY .. /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/
COPY .env .
WORKDIR /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5459
RUN go mod download
RUN go build -o ./.bin ./cmd/auth/main.go

ENTRYPOINT ["./.bin"]