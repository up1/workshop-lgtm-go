package main

import (
	"context"
	"demo/api"
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	initTracing()

	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := api.GetUser(c.Request.Context(), id)
		c.JSON(200, gin.H{
			"id":   id,
			"user": name,
		})
	})
	_ = r.Run(":8080")
}

func initTracing() {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}
	// Set sample rate to 1.0 for demonstration purposes.
	// In production, you might want to set a lower sample rate.

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// 1% is a common sample rate, but you can adjust it as needed.
		// trace.WithSampler(trace.TraceIDRatioBased(0.01)), // 1% sample rate
		trace.WithSampler(trace.AlwaysSample()), // Always sample for demo purposes
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}
