FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/server.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .

COPY server_config.yaml .

EXPOSE 8080
CMD ["/app/main"]