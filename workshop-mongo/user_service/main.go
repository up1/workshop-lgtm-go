package main

import (
	"context"
	"demo/api"
	"demo/db"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	initTracing()
	initMeter()
	initLogging()

	// Connect to MongoDB
	client, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	userService := &api.UserService{
		Client: client,
	}

	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		num, _ := strconv.Atoi(id)
		name := userService.GetUser(c.Request.Context(), num)
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

func initLogging() {
	logExporter, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)
	global.SetLoggerProvider(loggerProvider)
}
