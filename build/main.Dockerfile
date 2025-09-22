FROM golang:1.24.0-alpine AS builder

COPY .. /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/tree/dev/Elizaveta-Makeeva/
WORKDIR /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/tree/dev/Elizaveta-Makeeva/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go


FROM scratch AS runner

WORKDIR /kinopoisk-back/

COPY --from=builder /github.com/go-park-mail-ru/2025_2_DavaiDavaiDeploy/tree/dev/Elizaveta-Makeeva/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5458

ENTRYPOINT ["./.bin"]
