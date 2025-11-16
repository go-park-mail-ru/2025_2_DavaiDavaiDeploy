FROM golang:1.24.0-alpine AS builder

COPY .. /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/
WORKDIR /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go


FROM scratch AS runner

WORKDIR /kinopoisk-back/

COPY --from=builder /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /

COPY .env .
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip
ENV GOPROXY=https://proxy.golang.org,direct

EXPOSE 5458

ENTRYPOINT ["./.bin"]
