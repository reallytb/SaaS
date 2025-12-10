package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"SaaS/internal/dto"
	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

// CreateComment godoc
// @Summary Создать комментарий к задаче
// @Description Создает новый комментарий для указанной задачи. Требуются права редактора или владельца проекта.
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Param input body dto.CreateCommentRequest true "Данные комментария"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id}/comments [post]
// @Example request
//
//	{
//	    "content": "Это комментарий к задаче"
//	}
//
// @Example response 201
//
//	{
//	    "success": true,
//	    "message": "комментарий для задачи Реализовать API успешно создан"
//	}
//
// @Example response 403
//
//	{
//	    "success": false,
//	    "error": "недостаточно прав для создания комментария к задаче данного проекта"
//	}
//
// @Example response 404
//
//	{
//	    "success": false,
//	    "error": "задача не найдена"
//	}
func CreateComment(c *gin.Context) {
	user := utils.GetUser(c)
	taskId := c.Param("id")
	var project models.Project
	var task models.Task
	var comment models.Comment
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
				Error:   "недостаточно прав для создания комментария к задаче данного проекта",
			})
			return
		}
	}
	if err := c.Bind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	if strings.TrimSpace(comment.Content) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "комментарий не может быть пустым",
		})
		return
	}
	comment.TaskId = task.ID
	comment.UserId = user.ID
	result = initializers.DB.Create(&comment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка создания комментария",
		})
		return
	}
	message := fmt.Sprintf("комментарий для задачи %s успешно создан", task.Title)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// GetComments godoc
// @Summary Получить комментарии задачи
// @Description Возвращает все комментарии для указанной задачи. Требуются права просмотра проекта.
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /tasks/{id}/comments [get]
// @Example response 200 (с комментариями)
//
//	{
//	    "success": true,
//	    "data": [
//	        {
//	            "id": 1,
//	            "task_title": "Реализовать API",
//	            "user_name": "Иван Иванов",
//	            "content": "Это комментарий к задаче",
//	            "created_at": "2023-12-15 14:30:45"
//	        }
//	    ]
//	}
//
// @Example response 403
//
//	{
//	    "success": false,
//	    "error": "недостаточно прав для просмотра комментариев к задаче данного проекта"
//	}
//
// @Example response 200 (пустой список)
//
//	{
//	    "success": true,
//	    "message": "у задачи нет комментариев"
//	}
func GetComments(c *gin.Context) {
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
			Error:   "проект не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) == 0 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для просмотра комментариев задачи данного проекта",
			})
			return
		}
	}
	var comments []models.Comment
	var respComments []models.RespComment
	initializers.DB.Where("task_id = ?", task.ID).Find(&comments)
	if len(comments) == 0 {
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Success: true,
			Message: "у задачи нет комментариев",
		})
		return
	}
	for _, comment := range comments {
		var commentUser models.User
		initializers.DB.First(&commentUser, "ID = ?", comment.UserId)
		respComment := models.RespComment{
			ID:        comment.ID,
			TaskTitle: task.Title,
			UserName:  commentUser.Name,
			Content:   comment.Content,
			CreatedAt: utils.TimeFormat(comment.CreatedAt),
		}
		respComments = append(respComments, respComment)
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    respComments,
	})
}

// DeleteComment godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий по ID. Требуются права редактора или владельца проекта.
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комментария"
// @Success 200 {object} dto.SuccessResponse "Комментарий удален"
// @Failure 400 {object} dto.ErrorResponse "Некорректный запрос"
// @Failure 403 {object} dto.ErrorResponse "Недостаточно прав"
// @Failure 404 {object} dto.ErrorResponse "Комментарий не найден"
// @Router /api/comments/{id} [delete]
// @Example response 200
//
//	{
//	    "success": true,
//	    "message": "комментарий к задаче Реализовать API удален"
//	}
//
// @Example response 403
//
//	{
//	    "success": false,
//	    "error": "недостаточно прав для удаления комментария к задаче данного проекта"
//	}
//
// @Example response 404
//
//	{
//	    "success": false,
//	    "error": "комментарий не найден"
//	}
func DeleteComment(c *gin.Context) {
	user := utils.GetUser(c)
	commentId := c.Param("id")
	var project models.Project
	var task models.Task
	var comment models.Comment
	result := initializers.DB.First(&comment, "ID = ?", commentId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "комментарий не найден",
		})
		return
	}
	result = initializers.DB.First(&task, "ID = ?", comment.TaskId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "задача не найдена",
		})
		return
	}
	result = initializers.DB.First(&project, "ID = ?", task.ProjectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "проект не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для удаления комментария к задаче данного проекта",
			})
			return
		}
	}
	result = initializers.DB.Delete(&comment, "ID = ?", commentId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка удаления комментария"})
		return
	}
	message := fmt.Sprintf("комментарий к задаче %s удален", task.Title)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}
