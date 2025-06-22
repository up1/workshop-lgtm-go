package api

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var tracer = otel.Tracer("gin-server")

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	Client *mongo.Client
}

func (us *UserService) GetUser(ctx context.Context, id int) string {
	_, span := tracer.Start(ctx, "GetUser")
	defer span.End()

	var user User
	err := us.Client.Database("test").Collection("users").FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		countByStatus(ctx, "error")
		return "unknown"
	}
	countByStatus(ctx, "success")
	return user.Name
}

func countByStatus(ctx context.Context, status string) {
	commonAttrs := []attribute.KeyValue{
		attribute.String("status", status),
	}
	meter := otel.Meter("demo_meter")
	requestCountByStatus, err := meter.Int64Counter("request_get_user_count")
	if err != nil {
		log.Fatal(err)
	}
	requestCountByStatus.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
}
