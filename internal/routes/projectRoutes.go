package routes

import (
	"SaaS/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setUpProjectRoutes(r *gin.RouterGroup) {
	r.POST("/projects", handlers.CreateProject)
	r.GET("/projects", handlers.GetProjects)
	r.GET("/projects/:id", handlers.GetProject)
	r.PATCH("/projects/:id", handlers.EditProject)
	r.DELETE("/projects/:id", handlers.DeleteProject)
}
