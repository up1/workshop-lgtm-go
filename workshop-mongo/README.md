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

## Example with slog
```
{"Timestamp":"2025-06-22T17:49:17.803797885Z","ObservedTimestamp":"2025-06-22T17:49:17.803806135Z","Severity":9,"SeverityText":"INFO","Body":{"Type":"String","Value":"GetUser called"},"Attributes":[{"Key":"id","Value":{"Type":"Int64","Value":123}}],"TraceID":"ee60a7da339c1c1b2b43027909b054e0","SpanID":"b97b9d9de709ea13","TraceFlags":"01","Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"user_service"}},{"Key":"telemetry.sdk.language","Value":{"Type":"STRING","Value":"go"}},{"Key":"telemetry.sdk.name","Value":{"Type":"STRING","Value":"opentelemetry"}},{"Key":"telemetry.sdk.version","Value":{"Type":"STRING","Value":"1.36.0"}}],"Scope":{"Name":"gin-server","Version":"","SchemaURL":"","Attributes":{}},"DroppedAttributes":0}
```

## Resources
* [Otel receiver :: MongoDB ](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mongodbreceiver)
* [Otel slog](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/bridges/otelslog)