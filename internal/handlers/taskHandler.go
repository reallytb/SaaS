package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"SaaS/internal/dto"
	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

// CreateTask godoc
// @Summary Создать задачу
// @Description Создаёт новую запись в таблице задач. Требуются права редактора проекта.
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID проекта"
// @Param input body dto.CreateTaskRequest true "Название, описание, статус, приоритет и дата завершения задачи"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id}/tasks [post]
func CreateTask(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var task models.Task
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка получения тела запроса",
		})
		return
	}
	if task.Title == "" || len(task.Title) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "название задачи не может быть пустым",
		})
		return
	}
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект не найден",
		})
		return
	}
	task.ProjectId = project.ID
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для создания задачи в данном проекте",
			})
			return
		}
	}
	if task.Status != "todo" && task.Status != "in-progress" && task.Status != "done" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "укажите действительный статус задачи",
		})
		return
	}
	if task.Priority != "low" && task.Priority != "medium" && task.Priority != "high" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "укажите действительный приоритет задачи",
		})
		return
	}
	initializers.DB.Create(&task)
	message := fmt.Sprintf("задача %s создана", task.Title)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// GetTasks godoc
// @Summary Получить задачи проекта
// @Description Возвращает все задачи проекта. Требуются права просмотра проекта.
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID проекта"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id}/tasks [get]
func GetTasks(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) == 0 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для просмотра задач данного проекта",
			})
			return
		}
	}
	var tasks []models.Task
	var respTasks []models.Resptask
	initializers.DB.Where("Project_id = ?", project.ID).Find(&tasks)
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Success: true,
			Message: "у проекта нет задач",
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
			DueDate:     utils.TimeFormat(task.Due_date),
			CreatedAt:   utils.TimeFormat(task.CreatedAt),
			UpdatedAt:   utils.TimeFormat(task.UpdatedAt),
		}
		respTasks = append(respTasks, respTask)
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    respTasks,
	})
}

// GetTask godoc
// @Summary Получить задачу
// @Description Возвращает одну конкретную задачу. Требуются права просмотра проекта.
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id} [post]
func GetTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "задача не найдена",
		})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.ProjectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект задачи не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) == 0 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для просмотра задачи данного проекта",
			})
			return
		}
	}
	respTask := models.Resptask{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     utils.TimeFormat(task.Due_date),
		CreatedAt:   utils.TimeFormat(task.CreatedAt),
		UpdatedAt:   utils.TimeFormat(task.UpdatedAt),
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    respTask,
	})
}

// EditTask godoc
// @Summary Изменить задачу
// @Description Изменяет запись в таблице задач. Требуются права редактора проекта.
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param input body dto.CreateTaskRequest true "Название, описание, статус, приоритет и дата завершения задачи"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id} [patch]
func EditTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "задача не найдена",
		})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.ProjectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект задачи не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для редактирования задачи данного проекта",
			})
			return
		}
	}
	if c.Bind(&task) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	if task.Status != "todo" && task.Status != "in-progress" && task.Status != "done" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "укажите действительный статус задачи",
		})
		return
	}
	if task.Priority != "low" && task.Priority != "medium" && task.Priority != "high" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "укажите действительный приоритет задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("title", task.Title)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения названия задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("description", task.Description)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения описания задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("due_date", task.Due_date)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения срока задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("priority", task.Priority)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения приоритета задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("status", task.Status)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения статуса задачи",
		})
		return
	}
	result = initializers.DB.Model(&task).UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения времени обновления задачи",
		})
		return
	}
	message := fmt.Sprintf("задача %s успешно изменена", task.Title)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// DeleteTask godoc
// @Summary Удалить
// @Description Удаляет запись в таблице задач. Требуются права редактора проекта.
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	result := initializers.DB.First(&task, "ID = ?", taskId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "задача не найдена",
		})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.ProjectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект задачи не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав дла удаления задачи данного проекта",
			})
			return
		}
	}
	result = initializers.DB.Delete(&task, "ID = ?", task.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка удаления задачи",
		})
		return
	}
	message := fmt.Sprintf("задача %s успешно удалёна", task.Title)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}
