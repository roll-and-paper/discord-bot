FROM amd64/alpine as builder
RUN apk update && apk add ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY discord-bot /discord-bot

EXPOSE 8080

CMD ["/discord-bot", "start"]
