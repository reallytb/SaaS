package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"SaaS/internal/initializers"
	"SaaS/internal/models"
	"SaaS/pkg/utils"
)

func CreatePermission(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	var permission models.Permission
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	if user.ID != project.OwnerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "вы не являетесь владельцем данного проекта"})
		return
	}
	var body struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	if body.Role != "viewer" && body.Role != "editor" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "указана неверная роль"})
		return
	}
	permission.ProjectId = project.ID
	var permissionUser models.User
	result = initializers.DB.First(&permissionUser, "email = ?", body.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким email не найден"})
		return
	}
	var permissionCheck models.Permission
	permission.UserId = permissionUser.ID
	permission.Role = body.Role
	result = initializers.DB.First(&permissionCheck, "user_id = ? AND project_id = ?", permission.UserId, project.ID)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "у пользователя с таким email уже есть права на данный проект"})
		return
	}
	result = initializers.DB.Create(&permission)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка создания прав"})
		return
	}
	message := fmt.Sprintf("пользователь %s теперь имеет права %s для проекта %s", user.Name, permission.Role, project.Name)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func DeletePermission(c *gin.Context) {
	user := utils.GetUser(c)
	projectId := c.Param("id")
	userId := c.Param("user_id")
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	if user.ID != project.OwnerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "вы не являетесь владельцем данного проекта"})
		return
	}
	var permission models.Permission
	var userToDelete models.User
	result = initializers.DB.First(&userToDelete, "ID = ?", userId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "данный пользователь не найден"})
		return
	}
	result = initializers.DB.First(&permission, "user_id = ? AND project_id = ?", userToDelete.ID, project.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "права у данного пользователя на проект не найдены"})
		return
	}
	result = initializers.DB.Delete(&permission, "ID = ?", permission.ID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка удаления прав на проект"})
		return
	}
	message := fmt.Sprintf("права польозвателя %s на проект %s удалены", userToDelete.Name, project.Name)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
