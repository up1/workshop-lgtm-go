package api

import (
	"context"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gin-server")

func GetUser(ctx context.Context, id string) string {
	_, span := tracer.Start(ctx, "GetUser")
	defer span.End()
	if id == "123" {
		return "demo name"
	}
	return "unknown"
}
