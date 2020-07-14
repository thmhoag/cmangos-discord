FROM golang:1.14.4 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./bin/cmangos-discord ./cmd/cmangos-discord 
RUN chmod a+rx ./bin/cmangos-discord

# Output final image
FROM scratch as final

COPY --from=builder /app/bin/cmangos-discord /

VOLUME /config/cmangos-discord

ENTRYPOINT [ "/cmangos-discord" ]
