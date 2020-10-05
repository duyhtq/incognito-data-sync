FROM golang:1.12.9-alpine AS builder

ARG API_CONF_DATA
ARG GITHUB_TOKEN

WORKDIR /go/src/github.com/duyhtq/incognito-data-sync

RUN apk add --no-cache git dep

COPY . .
RUN dep ensure -v

RUN echo -n "$API_CONF_DATA" | base64 -d > config/conf.json

RUN go build -o server server.go
RUN go build -o cron-sync-ptrade ./extensions/cron_sync_ptrade/*.go
RUN go build -o cron-update-price ./extensions/cron_update_price/*.go

FROM alpine:3.11

WORKDIR /app

RUN apk add gcc libc-dev ca-certificates linux-headers

COPY --from=builder /go/src/github.com/duyhtq/incognito-data-sync/server .
COPY --from=builder /go/src/github.com/duyhtq/incognito-data-sync/cron-sync-ptrade .
COPY --from=builder /go/src/github.com/duyhtq/incognito-data-sync/cron-update-price .
COPY --from=builder /go/src/github.com/duyhtq/incognito-data-sync/config/conf.json ./config/conf.json

RUN chmod +x ./server

CMD ["./server"]
