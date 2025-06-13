FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 8080
CMD ["./order-service"]