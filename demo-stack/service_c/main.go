package main

import (
	"net/http"
	"shared"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Initialize OpenTelemetry tracing
	shared.InitTracing()

	// Create a new Gin router
	r := gin.New()

	// Middleware for OpenTelemetry tracing
	// This will automatically instrument incoming HTTP requests
	r.Use(otelgin.Middleware("my-server"))

	r.Use(gin.Recovery())

	// Get /products endpoint
	r.GET("/products", func(c *gin.Context) {
		// Create a new span for this request
		span1Ctx, span := shared.StartNewSpan(c.Request.Context(), "service-c", "GetProducts")
		defer span.End()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		// Create a span for the business logic
		_, businessLogicSpan := shared.StartNewSpan(span1Ctx, "service-c", "BusinessLogic")
		defer businessLogicSpan.End()

		c.JSON(http.StatusOK, gin.H{"products": []string{"Product 1", "Product 2"}})
	})

	r.Run(":8080")
}
