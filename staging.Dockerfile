FROM golang:1.13.12-alpine AS builder

RUN apk add gcc libc-dev ca-certificates linux-headers git

ARG GITHUB_TOKEN
ARG API_CONF_DATA

ENV GOPRIVATE github.com/inc-backend

WORKDIR /app

RUN git config --global --add url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/inc-backend".insteadOf "https://github.com/inc-backend"

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN echo -n "$API_CONF_DATA" | base64 -d > config/conf.json

RUN go build -o server server.go
RUN go build -o cron-sync-ptrade ./extensions/cron_sync_ptrade/*.go
RUN go build -o cron-sync-data ./extensions/cron_sync_data/*.go
RUN go build -o cron-update-price ./extensions/cron_update_price/*.go

FROM alpine:3.12.0

RUN apk --no-cache add ca-certificates gcc libc-dev linux-headers

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/cron-sync-ptrade .
COPY --from=builder /app/cron-sync-data .
COPY --from=builder /app/cron-update-price .
COPY --from=builder /app/config/conf.json ./config/conf.json

RUN chmod +x ./server

CMD ["./server"]
