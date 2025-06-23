FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN cd service_a && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o api service_a

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]