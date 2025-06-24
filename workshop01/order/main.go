package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	// Initialize OpenTelemetry tracing
	cleanup, err := setupTraceProvider(os.Getenv("OTEL_ENDPOINT"), os.Getenv("SERVICE_NAME"), "1.0.0")
	if err != nil {
		log.Fatalf("Failed to set up trace provider: %v", err)
	}
	defer cleanup()

	// Initialize Prometheus metrics
	router := gin.New()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	// Register OpenTelemetry middleware
	router.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME")))

	// Initialize Slog logger
	// with span ID and trace ID enabled
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config := sloggin.Config{
		WithRequestID: true,
		WithSpanID:    true,
		WithTraceID:   true,
	}
	router.Use(sloggin.NewWithConfig(logger, config))

	router.Use(gin.Recovery())

	// Example route
	router.GET("/ping", func(c *gin.Context) {
		sloggin.AddCustomAttributes(c, slog.String("ping", "pong"))
		sloggin.AddCustomAttributes(c, slog.String("ping2", "pong2"))

		// Log a message
		slog.Info("Received ping request 2", slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path))
		slog.Info("Received ping request 3", slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path))

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}

func setupTraceProvider(endpoint string, serviceName string, serviceVersion string) (func(), error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	cleanup := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}
	return cleanup, nil
}
