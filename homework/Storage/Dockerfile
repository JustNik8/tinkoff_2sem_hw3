FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main ./cmd/storage.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .

COPY server_config.yaml .

CMD ["/app/main"]