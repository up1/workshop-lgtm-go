package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

type UserHandler struct {
	Ch *amqp091.Channel
}

func (h *UserHandler) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {

		createMetric(c)

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Publish the user creation event to a message queue (rabbitmq)
		PublishUserCreationEvent(User{
			Name:  req.Name,
			Email: req.Email,
		}, c.Request.Context(), h.Ch)

		c.JSON(http.StatusCreated, CreateUserResponse{
			Message: "User created successfully",
			User: User{
				ID:    1, // Mock ID
				Name:  req.Name,
				Email: req.Email,
			},
		})
	}
}

func createMetric(c *gin.Context) {
	meter := otel.Meter("demo_meter")
	commonAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
	}
	// request_count_total
	requestCount, err := meter.Int64Counter("service_a_request_count_total")
	if err != nil {
		log.Fatal(err)
	}
	requestCount.Add(c.Request.Context(), 1, metric.WithAttributes(commonAttrs...))
}
