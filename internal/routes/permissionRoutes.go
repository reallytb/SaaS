package routes

import (
	"SaaS/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setUpPermissionRoutes(r *gin.RouterGroup) {
	r.POST("/projects/:id/share", handlers.CreatePermission)
	r.DELETE("/projects/:id/share/:user_id", handlers.DeletePermission)
}
