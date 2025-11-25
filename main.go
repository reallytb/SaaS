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

	r.POST("/projects/:id/tasks", middleware.AuthCheck, handlers.CreateTask)
	r.GET("/projects/:id/tasks", middleware.AuthCheck, handlers.GetTasks)
	r.GET("/tasks/:id", middleware.AuthCheck, handlers.GetTask)
	r.PATCH("/tasks/:id", middleware.AuthCheck, handlers.EditTask)
	r.DELETE("/tasks/:id", middleware.AuthCheck, handlers.DeleteTask)

	r.POST("/tasks/:id/comments", middleware.AuthCheck, handlers.CreateComment)
	r.GET("/tasks/:id/comments", middleware.AuthCheck, handlers.GetComments)
	r.DELETE("/comments/:id", middleware.AuthCheck, handlers.DeleteComment)

	r.POST("/projects/:id/share", middleware.AuthCheck, handlers.CreatePermission)
	r.DELETE("/projects/:id/share/:user_id", middleware.AuthCheck, handlers.DeletePermission)

	r.Run()
}
