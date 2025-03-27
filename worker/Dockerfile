FROM golang:1.23.1 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
RUN chmod +x ./app  # Даем права на выполнение
CMD ["./app"]