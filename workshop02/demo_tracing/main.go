package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	rdb *redis.Client
)

func initTracing() {
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}
	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exporter))
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func initMeter() {
	metricExporter, err := otlpmetrichttp.New(context.Background())
	if err != nil {
		panic(err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, rootSpan := otel.Tracer("redis01").Start(r.Context(), "handleRequest")
	defer rootSpan.End()

	meter := otel.Meter("demo_meter")
	commonAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
	}
	// request_count_total
	requestCount, err := meter.Int64Counter("request_count_total")
	if err != nil {
		log.Fatal(err)
	}
	requestCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))

	cmd := rdb.Incr(ctx, "counter")
	if err := cmd.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(strconv.FormatInt(cmd.Val(), 10)))

}

func main() {
	rdb = redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_SERVER")})
	h := http.Handler(http.HandlerFunc(handler))
	if os.Getenv("ENABLE_OTEL") != "" {
		log.Println("enabling opentelemetry")
		initTracing()
		initMeter()
		h = otelhttp.NewHandler(http.HandlerFunc(handler), "GET /")
	}
	log.Fatal(http.ListenAndServe(":8080", h))
}
