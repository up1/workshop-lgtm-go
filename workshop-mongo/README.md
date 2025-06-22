# Workshop with MongoDB
* API with Go + Gin
* MongoDB


## Develop REST API with Go + Gin
* Tracing
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