package user

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func PublishUserCreationEvent(user User, ctx context.Context, ch *amqp091.Channel) {
	// Start a new span for the message publishing operation
	span1 := startNewSpanWithAttributes(ctx, "SerializedUserCreationEvent",
		attribute.Int("user_id", user.ID),
		attribute.String("user_name", user.Name),
		attribute.String("user_email", user.Email),
		attribute.String("operation", "serialize"),
	)
	defer span1.End()

	span2 := startNewSpanWithAttributes(ctx, "PublishUserCreationEvent",
		attribute.String("exchange_name", "users"),
		attribute.String("operation", "publish"),
	)
	defer span2.End()

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

func startNewSpan(ctx context.Context, operationName string) trace.Span {
	_, span := otel.Tracer("service-a").Start(ctx, operationName)
	return span
}

func startNewSpanWithAttributes(ctx context.Context, operationName string, attributes ...attribute.KeyValue) trace.Span {
	span := startNewSpan(ctx, operationName)
	span.SetAttributes(attributes...)
	return span
}
