package routes

import (
	"SaaS/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setUpCommentRoutes(r *gin.RouterGroup) {
	r.POST("/tasks/:id/comments", handlers.CreateComment)
	r.GET("/tasks/:id/comments", handlers.GetComments)
	r.DELETE("/comments/:id", handlers.DeleteComment)
}
