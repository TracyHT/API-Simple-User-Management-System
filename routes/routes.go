package routes

import (
	"github.com/gin-gonic/gin"

	"main.go/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.GET("/users", handlers.ListUsers)
}
