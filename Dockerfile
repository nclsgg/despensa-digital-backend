FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]