package routes

import (
	"SaaS/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setUpAuthRoutes(r *gin.Engine) {
	r.POST("/signup", handlers.SignUp)
	r.POST("/signin", handlers.SignIn)
}
