FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o aggregator ./cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/aggregator .
EXPOSE 8080
CMD ["./aggregator"]