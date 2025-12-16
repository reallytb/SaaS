package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"SaaS/internal/dto"
	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

// CreatePermission godoc
// @Summary Выдать права другому пользователю
// @Description Создаёт новую запись в таблице прав. Требуются права владельца проекта.
// @Tags permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID проекта"
// @Param input body dto.CreatePermissionRequest true "Email и роль для выдачи роли"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id}/share [post]
// @Example request
//
//	{
//	   "email": "example@example.com"
//	   "role": "editor"
//
// }
// @Example response 201
//
//	{
//	    "success": true,
//	    "message": "пользователь Иван теперь имеет права viewer для проекта Написать API"
//	}
//
// @Example response 403
//
//	{
//	    "success": false,
//	    "error": "вы не являетесь владельцем данного проекта"
//	}
//
// @Example response 400 (указана неверная роль)
//
//	{
//	    "success": true,
//	    "message": "указана неверная роль"
//	}
func CreatePermission(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	var permission models.Permission
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "проект не найден",
		})
		return
	}
	if user.ID != project.OwnerId {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Error:   "вы не являетесь владельцем данного проекта",
		})
		return
	}
	var body struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка чтения тела запроса",
		})
		return
	}
	if body.Role != "viewer" && body.Role != "editor" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "указана неверная роль",
		})
		return
	}
	permission.ProjectId = project.ID
	var permissionUser models.User
	result = initializers.DB.First(&permissionUser, "email = ?", body.Email)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "пользователь с таким email не найден",
		})
		return
	}
	var permissionCheck models.Permission
	permission.UserId = permissionUser.ID
	permission.Role = body.Role
	result = initializers.DB.First(&permissionCheck, "user_id = ? AND project_id = ?", permission.UserId, project.ID)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "у пользователя с таким email уже есть права на данный проект",
		})
		return
	}
	result = initializers.DB.Create(&permission)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка создания прав",
		})
		return
	}
	message := fmt.Sprintf("пользователь %s теперь имеет права %s для проекта %s", user.Name, permission.Role, project.Name)
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}

// DeletePermission godoc
// @Summary Удалить права у пользователя
// @Description Удаляет запись в таблице прав. Требуются права владельца проекта.
// @Tags permission
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID проекта"
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /projects/{id}/share/{user_id} [delete]
// @Example response 200
//
//	{
//	    "success": true,
//	    "message": "права пользователя Иван на проект Написать API удалены"
//	}
//
// @Example response 403
//
//	{
//	    "success": false,
//	    "error": "вы не являетесь владельцем данного проекта"
//	}
func DeletePermission(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	userId := c.Param("user_id")
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
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Error:   "вы не являетесь владельцем данного проекта",
		})
		return
	}
	var permission models.Permission
	var userToDelete models.User
	result = initializers.DB.First(&userToDelete, "ID = ?", userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "пользователь не найден",
		})
		return
	}
	result = initializers.DB.First(&permission, "user_id = ? AND project_id = ?", userToDelete.ID, project.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "права на данный проект у пользователя не найдены",
		})
		return
	}
	result = initializers.DB.Delete(&permission, "ID = ?", permission.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "ошибка удаления прав на проект",
		})
		return
	}
	message := fmt.Sprintf("права польозвателя %s на проект %s удалены", userToDelete.Name, project.Name)
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: message,
	})
}
