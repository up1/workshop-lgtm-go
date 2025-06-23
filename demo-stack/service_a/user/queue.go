package user

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func PublishUserCreationEvent(user User, ctx context.Context, ch *amqp091.Channel) {
	// Start a new span for the message publishing operation
	_, span := otel.Tracer("service-a").Start(ctx, "publish-rabbitmq-message")
	// Add channel information to the span
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "exchange.name",
			Value: attribute.StringValue("users"),
		},
	)
	defer span.End()

	// Inject the context into the message headers for tracing
	headers := make(map[string]string)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator.Inject(ctx, propagation.MapCarrier(headers))

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user: %v", err)
		return
	}
	amqpHeaders := make(amqp091.Table)
	for k, v := range headers {
		amqpHeaders[k] = v
	}

	// Publish the user creation event to RabbitMQ
	err = ch.PublishWithContext(ctx,
		"users", // exchange
		"",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        userJSON,
			Headers:     amqpHeaders,
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Println("User creation event published successfully")
	}
}
