package routes

import (
	"github.com/gin-gonic/gin"

	"SaaS/internal/handlers"
)

func setUpUserRoutes(r *gin.RouterGroup) {
	r.GET("/me", handlers.Me)
	r.POST("/signout", handlers.SignOut)
}
