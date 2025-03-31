FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build .

FROM alpine:latest AS runner

WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./main"]