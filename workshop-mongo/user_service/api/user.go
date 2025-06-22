package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
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
		return "unknown"
	}
	return user.Name
}
