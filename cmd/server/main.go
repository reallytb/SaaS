package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "SaaS/docs"
	"SaaS/internal/initializers"
	"SaaS/internal/routes"
)

// @title SaaS
// @version 1.0
// @description API для управления задачами, проектами и комментариями
// @termsOfService http://swagger.io/terms/

// @contact.name Support Team
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer {token}" (без кавычек)
func main() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.SyncDB()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.SetUpRoutes(r)

	r.Run()
}
