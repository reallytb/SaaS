package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

func CreateTask(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var task models.Task
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка получения тела запроса"})
		return
	}
	if task.Title == "" || len(task.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "название задачи не может быть пустым"})
		return
	}
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	task.Project_id = project.ID
	if user.ID != project.Owner_id {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для создания задачи в данном проекте"})
			return
		}
	}
	initializers.DB.Create(&task)
	message := fmt.Sprintf("задача %s создана", task.Title)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func GetTasks(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	if user.ID != project.Owner_id {
		if utils.PermissionCheck(user, project) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для просмотра задач данного проекта"})
			return
		}
	}
	var tasks []models.Task
	var respTasks []models.Resptask
	initializers.DB.Where("Project_id = ?", project.ID).Find(&tasks)
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "у проекта нет задач",
		})
		return
	}
	for _, task := range tasks {
		respTask := models.Resptask{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			Due_date:    utils.TimeFormat(task.Due_date),
			Created_at:  utils.TimeFormat(task.CreatedAt),
			Updated_at:  utils.TimeFormat(task.UpdatedAt),
		}
		respTasks = append(respTasks, respTask)
	}
	c.JSON(http.StatusOK, gin.H{
		"tasks": respTasks,
	})
}

func GetTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "задача не найдена"})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.Project_id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект задачи не найден"})
		return
	}
	if user.ID != project.Owner_id {
		if utils.PermissionCheck(user, project) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для просмотра задачи данного проекта"})
			return
		}
	}
	respTask := models.Resptask{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Due_date:    utils.TimeFormat(task.Due_date),
		Created_at:  utils.TimeFormat(task.CreatedAt),
		Updated_at:  utils.TimeFormat(task.UpdatedAt),
	}
	c.JSON(http.StatusOK, respTask)
}

func EditTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "задача не найдена"})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.Project_id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект задачи не найден"})
		return
	}
	if user.ID != project.Owner_id {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для редактирования задачи данного проекта"})
			return
		}
	}
	if c.Bind(&task) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("title", task.Title)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения названия задачи"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("description", task.Description)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения описания задачи"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("due_date", task.Due_date)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения срока задачи"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("priority", task.Priority)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения приоритета задачи"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("status", task.Status)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения статуса задачи"})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения времени изменения задачи"})
		return
	}
	message := fmt.Sprintf("задача %s успешно изменена", task.Title)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func DeleteTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "задача не найдена"})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.Project_id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект задачи не найден"})
		return
	}
	if user.ID != project.Owner_id {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для удаления задачи данного проекта"})
			return
		}
	}
	result = initializers.DB.Delete(&task, "ID = ?", task.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка удаления задачи"})
		return
	}
	message := fmt.Sprintf("задача %s успешно удалёна", task.Title)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
