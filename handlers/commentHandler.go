package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"SaaS/initializers"
	"SaaS/models"
	"SaaS/services"
)

func CreateComment(c *gin.Context) {
	user := services.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	var comment models.Comment
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "задача не найдена"})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.Project_id)
	if user.ID != project.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для просмотра задачи данного проекта"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект задачи не найден"})
		return
	}
	if c.Bind(&comment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	comment.Task_id = task.ID
	comment.User_id = user.ID
	result = initializers.DB.Create(&comment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка создания комментария"})
		return
	}
	message := fmt.Sprintf("комментарий для задачи %s успешно создан", task.Title)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}
