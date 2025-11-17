package services

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"SaaS/models"
)

func GetUser(c *gin.Context) models.User {
	userGet, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователя не существует"})
		return models.User{}
	}
	user, ok := userGet.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный тип пользователя"})
		return models.User{}
	}
	return user
}
