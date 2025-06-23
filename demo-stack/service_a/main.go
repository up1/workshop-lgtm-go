package main

import (
	"net/http"
	"service_a/user"
	"shared"

	"github.com/gin-gonic/gin"
)

func main() {
	shared.InitTracing()

	r := gin.New()
	r.Use(gin.Recovery())

	// Create a new user
	r.POST("/user", func(c *gin.Context) {
		var newUser = user.CreateUserRequest{}
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.CreateUser(newUser)(c)
	})

	r.Run(":8080")
}
