FROM golang:1.24 AS builder

WORKDIR /build

COPY . .

RUN go mod tidy && \
    go build -o api .

FROM debian:stable-slim

WORKDIR /

# Install CA certificates
RUN apt update && apt add ca-certificates

COPY --from=builder /build/api .

ENV GIN_MODE=release

CMD ["/api"]