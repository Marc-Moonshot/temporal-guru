FROM golang:1.23-alpine AS builder

WORKDIR /sync-service

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main

FROM alpine:latest

WORKDIR /sync-service

COPY --from=builder /app/main .

EXPOSE 9000

CMD ["./main"]
