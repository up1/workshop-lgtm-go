# Workshop with MongoDB
* API with Go + Gin
* MongoDB

## Start MongoDB
```
$docker compose up -d mongo
$docker compose ps
```

## Develop REST API with Go + Gin
* Use go-mongo-driver v2
* Tracing with Gin and MongoDB
  * [Instrumentation with Gin](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation)


Build and start
```
$docker compose build user_service
$docker compose up -d user_service
$docker compose ps
```

List of URLs
* http://localhost:8080/users/123
* http://localhost:8080/users/1

Error message about Otel collector
```
traces export: Post "http://otel-collector:4318/v1/traces": dial tcp: lookup otel-collector on 127.0.0.11:53: no such host
```

## Start Otel collector server
* Metric with Prometheus
* Tracing with Jaeger

```
$docker compose up -d jaeger
$docker compose up -d prometheus
$docker compose up -d otel-collector

$docker compose ps
```

Go to Jaeger UI
* http://localhost:16686/search

Go to Prometheus UI
* http://localhost:9090

Call APIs again
* http://localhost:8080/users/123
* http://localhost:8080/users/1