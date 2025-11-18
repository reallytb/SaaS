package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"SaaS/initializers"
	"SaaS/models"
	"SaaS/services"
)

func CreateProject(c *gin.Context) {
	var project models.Project
	if c.Bind(&project) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	user := services.GetUser(c)
	project.Owner_id = user.ID
	result := initializers.DB.Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка создания проекта"})
		return
	}
	message := fmt.Sprintf("проект %s успешно создан", project.Name)
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

func GetProjects(c *gin.Context) {
	user := services.GetUser(c)
	var projects []models.Project
	var resProjects []models.Resproject
	initializers.DB.Where("Owner_id = ?", user.ID).Find(&projects)
	for _, project := range projects {
		resproject := models.Resproject{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Created_at:  services.TimeFormat(project.CreatedAt),
			Updated_at:  services.TimeFormat(project.UpdatedAt),
		}
		resProjects = append(resProjects, resproject)
	}
	fmt.Println(projects)
	if projects == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "у вас нет проектов",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"projects": resProjects,
	})
}

func GetProject(c *gin.Context) {
	user := services.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	if user.ID != project.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для просмотра данного проекта"})
		return
	}
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	resproject := models.Resproject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Created_at:  services.TimeFormat(project.CreatedAt),
		Updated_at:  services.TimeFormat(project.UpdatedAt),
	}
	c.JSON(http.StatusOK, resproject)
}

func EditProject(c *gin.Context) {
	user := services.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	if user.ID != project.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для редактирования данного проекта"})
		return
	}
	if c.Bind(&project) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения тела запроса"})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("name", project.Name)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения названия проекта"})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("description", project.Description)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка изменения описания проекта"})
		return
	}
	result = initializers.DB.Model(&project).UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка даты изменения проекта"})
		return
	}
	resproject := models.Resproject{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Created_at:  services.TimeFormat(project.CreatedAt),
		Updated_at:  services.TimeFormat(project.UpdatedAt),
	}
	c.JSON(http.StatusOK, resproject)
}

func DeleteProject(c *gin.Context) {
	user := services.GetUser(c)
	projectId := c.Param("id")
	var project models.Project
	result := initializers.DB.First(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "проект не найден"})
		return
	}
	if user.ID != project.Owner_id {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет прав для удаления данного проекта"})
		return
	}
	result = initializers.DB.Delete(&project, "ID = ?", projectId)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка удаления проекта"})
		return
	}
	message := fmt.Sprintf("проект %s успешно удалён", project.Name)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
