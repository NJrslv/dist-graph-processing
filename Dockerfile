# Build
FROM golang:1.21.4 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app/main src/main/test0.go

# Run
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app .
CMD ["app/main"]