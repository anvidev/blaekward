FROM golang:1.23.4 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o blaekward-server ./cmd/server/*.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/blaekward-server .
EXPOSE 9090 9091
ENTRYPOINT ["./blaekward-server"]
