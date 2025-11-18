package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"SaaS/initializers"
	"SaaS/models"
	"SaaS/services"
)

func CreateTask(c *gin.Context) {
	user := services.GetUser(c)
	projectId := c.Param("id")
	var task models.Task
	if err := c.Bind(&task); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка получения тела запроса"})
		return
	}
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка получения тела запроса"})
		return
	}
	task.Project_id = project.ID
	if project.Owner_id != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для создания задачи в данном проекте"})
		return
	}
	initializers.DB.Create(&task)
	message := fmt.Sprintf("задача %s создана", task.Title)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})

}
