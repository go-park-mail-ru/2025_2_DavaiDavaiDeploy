FROM golang:1.24.0-alpine AS builder

COPY .. /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/
WORKDIR /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/search/main.go


FROM scratch AS runner

WORKDIR /kinopoisk-back/

COPY --from=builder /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /

COPY .env .
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5462

ENTRYPOINT ["./.bin"]
