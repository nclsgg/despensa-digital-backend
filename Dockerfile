FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go mod download
RUN go build -o main ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod +x main

ENTRYPOINT ["./main"]
