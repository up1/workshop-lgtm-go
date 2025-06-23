# Workshop observability
* Service to service
  * HTTP
  * RabbitMQ
  * MongoDB
* Golang
  * Go workspace
  * Go modules

## Start RabbitMQ
```
$docker compose up -d rabbitmq
$docker compose ps
```

Access to managment UI
* http://localhost:15672/
  * username=guest
  * password=guest

## Build service-a
```
$docker compose build service_a
$docker compose up -d service_a
$docker compose ps
```

List of APIs
* POST /user
```
{
    "name": "somkiat",
    "email": "somkiat@xxx.com"
}
```

### Error message of Service-a
```
traces export: Post "http://otel-collector:4318/v1/traces": dial tcp: lookup otel-collector on 127.0.0.11:53: no such host
```

## Start Otel Collector server
```
$docker compose up -d jaeger
$docker compose up -d otel-collector
$docker compose ps
```

Go to Jaeger Server
* http://localhost:16686/search

