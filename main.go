package main

import (
	"github.com/gin-gonic/gin"

	"SaaS/handlers"
	"SaaS/initializers"
	"SaaS/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.POST("/signup", handlers.SignUp)
	r.POST("/signin", handlers.SignIn)
	r.GET("/me", middleware.AuthCheck, handlers.Me)
	r.POST("/signout", middleware.AuthCheck, handlers.SignOut)
	r.POST("/projects", middleware.AuthCheck, handlers.CreateProject)
	r.GET("/projects", middleware.AuthCheck, handlers.GetProjects)
	r.GET("/projects/:id", middleware.AuthCheck, handlers.GetProject)
	r.PATCH("/projects/:id", middleware.AuthCheck, handlers.EditProject)
	r.DELETE("/projects/:id", middleware.AuthCheck, handlers.DeleteProject)

	r.Run()
}
