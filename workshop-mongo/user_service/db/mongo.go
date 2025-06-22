package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
)

func Connect() (*mongo.Client, error) {
	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI("mongodb://localhost:27017/test?connect=direct")
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	// Ping the server to verify connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return client, nil
}
