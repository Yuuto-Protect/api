FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o api .

FROM debian:stable-slim

WORKDIR /

COPY --from=builder /app/api .

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/api"]