FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/avito

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
