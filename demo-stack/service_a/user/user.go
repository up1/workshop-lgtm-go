package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func CreateUser(req CreateUserRequest) func(c *gin.Context) {
	return func(c *gin.Context) {
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
