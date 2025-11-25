package main

import (
	"github.com/gin-gonic/gin"

	"SaaS/internal/initializers"
	"SaaS/internal/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()

	routes.SetUpRoutes(r)

	r.Run()
}
