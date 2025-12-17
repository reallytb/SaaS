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

// CreateProject godoc
// @Summary Создать проект
// @Description Создаёт новую запись в таблице проектов.
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateProjectRequest true "название и описание проекта"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /projects [post]
// @Example request
//
//	{
//	  "name": "Написать API",
//	  "description": "описание проекта"
//
//	}
//
// @Example response 201
//
//	{
//	    "success": true,
//	    "message": "проект Написать API успешно создан"
//	}
//
// @Example response 400
//
//	{
//	    "success": false,
//	    "error": "название проекта не может быть пустым"
//	}
func CreateProject(c *gin.Context) {
	var project models.Project
	if c.Bind(&project) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	if project.Name == "" || len(project.Name) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "название проекта не может быть пустым",
		})
		return
	}
	user := utils.GetUser(c)
	project.OwnerId = user.ID
	result := initializers.DB.Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка создания проекта",
		})
		return
	}
	message := fmt.Sprintf("проект %s успешно создан", project.Name)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// GetProjects godoc
// @Summary Получить список проектов
// @Description Возвращает список проектов авторизованного пользователя и их описание. Требуются права просмотра проекта.
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Router /projects [get]
// @Example response 200
//
//			{
//			    "success": true,
//			    "data": [
//			        {
//			            "id": 1,
//			            "name": "Написать API",
//			            "Description": "Описание проекта",
//			            "created_at": "2025-12-15 14:30:45",
//	                 "updated_at": "2025-12-15 14:30:45"
//			        }
//			    ]
//		 }
//
// @Example response 200 (нет проектов)
//
//	{
//		"Success": true,
//	    "Message": "У вас нет проектов"
//	}
func GetProjects(c *gin.Context) {
	user := utils.GetUser(c)
	var projects []models.Project
	var resProjects []models.Resproject
	initializers.DB.Where("Owner_id = ?", user.ID).Find(&projects)
	if len(projects) == 0 {
		c.JSON(http.StatusOK, dto.SuccessResponse{
			Success: true,
			Message: "У вас нет проектов",
		})
		return
	}
	for _, project := range projects {
		resproject := models.Resproject{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   utils.TimeFormat(project.CreatedAt),
			UpdatedAt:   utils.TimeFormat(project.UpdatedAt),
		}
		resProjects = append(resProjects, resproject)
	}
	fmt.Println(projects)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    resProjects,
	})
}

// GetProject godoc
// @Summary Получить конкретный проект
// @Description Возвращает проект и его описание. Требуются права просмотра проекта.
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id проекта"
// @Success 200 {object} dto.SuccessResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id} [get]
// @Example response 200
//
//			{
//			    "success": true,
//			    "data": [
//			        {
//			            "id": 1,
//			            "name": "Написать API",
//			            "Description": "Описание проекта",
//			            "created_at": "2025-12-15 14:30:45",
//	                 "updated_at": "2025-12-15 14:30:45"
//			        }
//			    ]
//		 }
//
// @Example response 404
//
//	{
//		"Success": false,
//	    "Error": "проект не найден"
//	}
func GetProject(c *gin.Context) {
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
				Error:   "недостаточно прав для просмотра данного проекта",
			})
			return
		}
	}
	resproject := models.Resproject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   utils.TimeFormat(project.CreatedAt),
		UpdatedAt:   utils.TimeFormat(project.UpdatedAt),
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    resproject,
	})
}

// EditProject godoc
// @Summary Изменить проект
// @Description Изменяет данные проекта на данные из тела запроса. Требуются права редактора проекта.
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id проекта"
// @Param input body dto.CreateProjectRequest true "название и описание проекта для редактирования"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id} [patch]
// @Example response 200
//
//				{
//				    "success": true,
//	             "message": "проект успешно изменен",
//				    "data": [
//				        {
//				            "id": 1,
//				            "name": "Написать API",
//				            "Description": "Описание проекта",
//				            "created_at": "2025-12-15 14:30:45",
//		                 "updated_at": "2025-12-17 8:16:00"
//				        }
//				    ]
//			 }
//
// @Example response 404
//
//	{
//		"Success": false,
//	    "Error": "проект не найден"
//	}
func EditProject(c *gin.Context) {
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
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для редактирования данного проекта",
			})
			return
		}
	}
	if c.Bind(&project) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("name", project.Name)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения названия проекта",
		})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("description", project.Description)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения описания проекта",
		})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка изменения даты обновления проекта",
		})
		return
	}
	resproject := models.Resproject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   utils.TimeFormat(project.CreatedAt),
		UpdatedAt:   utils.TimeFormat(project.UpdatedAt),
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "проект успешно изменен",
		Data:    resproject,
	})
}

// DeleteProject godoc
// @Summary Удалить проект
// @Description Удаляет запись из таблицы проектов. Требуются права редактора проекта.
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id проекта"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id} [delete]
// @Example response 200
//
//				{
//				  "success": true,
//	             "message": "проект Написать API успешно удалён"
//
// @Example response 404
//
//	{
//		"Success": false,
//	    "Error": "проект не найден"
//	}
func DeleteProject(c *gin.Context) {
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
		if utils.PermissionCheck(user, project) != 2 {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "недостаточно прав для удаления данного проекта",
			})
			return
		}
	}
	result = initializers.DB.Delete(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка удаления проекта",
		})
		return
	}
	message := fmt.Sprintf("проект %s успешно удалён", project.Name)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}
