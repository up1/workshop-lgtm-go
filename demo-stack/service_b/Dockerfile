FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN cd service_b && go mod tidy 
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer service_b

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]