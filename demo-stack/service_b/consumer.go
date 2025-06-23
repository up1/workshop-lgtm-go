package main

import (
	"context"
	"log"
	"os"
	"shared"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	shared.InitTracing()

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"users",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"users", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Headers
			headers := make(map[string]string)
			for k, v := range d.Headers {
				if str, ok := v.(string); ok {
					headers[k] = str
				}
			}
			processData(headers, d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func processData(headers map[string]string, body []byte) {

	// Create a new context with the extracted headers
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	parentCtx := propagator.Extract(context.Background(), propagation.MapCarrier(headers))
	ctx1, span1 := otel.Tracer("service_b").Start(parentCtx, "consume-rabbitmq-message")
	defer span1.End()

	log.Printf(" [x] Received a message with headers: %v", headers)
	log.Printf(" [x] %s", body)

	// Simulate processing
	log.Printf(" [x] Processing message: %s", body)
	time.Sleep(2 * time.Second)
	log.Printf(" [x] Finished processing message: %s", body)

	// Create a new span for processing the message
	_, span2 := otel.Tracer("service_b").Start(ctx1, "process-rabbitmq-message")
	defer span2.End()

}
