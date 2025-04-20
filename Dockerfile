FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.* .
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o blaekward-server ./cmd/server/*.go

FROM alpine:3.18 AS deployer
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/blaekward-server .
ENTRYPOINT ["./blaekward-server"]
