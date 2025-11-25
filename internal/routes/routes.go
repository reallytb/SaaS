package routes

import (
	"github.com/gin-gonic/gin"

	"SaaS/internal/middleware"
)

func SetUpRoutes(r *gin.Engine) {
	setUpAuthRoutes(r)

	protected := r.Group("/")
	protected.Use(middleware.AuthCheck)

	setUpCommentRoutes(protected)
	setUpPermissionRoutes(protected)
	setUpProjectRoutes(protected)
	setUpTaskRoutes(protected)
	setUpUserRoutes(protected)

}
