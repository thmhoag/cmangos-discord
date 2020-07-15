FROM golang:1.14.4-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o ./bin/cmangos-discord ./cmd/cmangos-discord
RUN chmod a+rx ./bin/cmangos-discord

# Output final image
FROM scratch as final

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/cmangos-discord /cmangos-discord

VOLUME /config/cmangos-discord

ENTRYPOINT [ "/cmangos-discord" ]
