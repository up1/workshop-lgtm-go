# Workshop with LGTM and Otel collector
* Open Telemetry
* [Otel collector](https://opentelemetry.io/docs/collector/)
  * The OpenTelemetry Collector offers a vendor-agnostic implementation of how to receive process and export telemetry data.
* LGTM Stack

### Start Redis server
```
$docker compose up -d redis
$docker compose ps
```

## Tracing with Jaeger
```
$docker compose up -d jaeger
$docker compose ps
```

Go to Jaeger UI
* http://localhost:16686/search

## Metric with Prometheus
```
$docker compose up -d prometheus
$docker compose ps
```

Go to Prometheus UI
* http://localhost:9090

## Install Otel Collector
* [Prometheus exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/prometheusexporter)
```
$docker compose up -d otel-collector
$docker compose ps
```

## Start demo for tracing
* Go + gin
* Redis
  * https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/redisreceiver/documentation.md


### Build and run demo
```
$docker compose build demo_tracing
$docker compose up -d demo_tracing
$docker compose ps
```
List of URLs
* http://localhost:8080/
