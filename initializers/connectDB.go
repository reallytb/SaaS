package initializers

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"SaaS/models"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("SaaS.db"), &gorm.Config{})

	if err != nil {
		panic("ошибка открытия базы данных")
	}
}

func SyncDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Project{})
	DB.AutoMigrate(&models.Task{})
	DB.AutoMigrate(&models.Comment{})
}
