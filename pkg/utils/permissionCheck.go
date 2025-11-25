package utils

import (
	"SaaS/internal/initializers"
	"SaaS/internal/models"
)

func PermissionCheck(user models.User, project models.Project) int {
	var permission models.Permission
	result := initializers.DB.First(&permission, "user_id = ? AND project_id = ?", user.ID, project.ID)
	if result.Error != nil {
		return 0
	}
	if permission.Role == "viewer" {
		return 1
	}
	if permission.Role == "editor" {
		return 2
	}
	return 0
}
