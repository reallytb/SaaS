package routes

import (
	"github.com/gin-gonic/gin"

	"SaaS/internal/handlers"
)

func setUpTaskRoutes(r *gin.RouterGroup) {
	r.POST("/projects/:id/tasks", handlers.CreateTask)
	r.GET("/projects/:id/tasks", handlers.GetTasks)
	r.GET("/tasks/:id", handlers.GetTask)
	r.PATCH("/tasks/:id", handlers.EditTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)
}
