FROM golang:1.19.4-alpine as builder
WORKDIR /data
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM caddy:2.6.4
WORKDIR /etc/caddy/
COPY --from=builder /data/CaddyfileTemplate .
COPY --from=builder /data/caddy-for-proxy .
COPY --from=builder /data/docker-entrypoint.sh .
RUN apk add --no-cache bash && chmod +x ./docker-entrypoint.sh
STOPSIGNAL SIGTERM
CMD ["./docker-entrypoint.sh"]
