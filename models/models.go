package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `json:"email,omitempty" gorm:"unique"`
	Password_hash string `json:"password_hash,omitempty"`
	Name          string `json:"name,omitempty"`
}

type Project struct {
	gorm.Model
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Owner_id    uint   `json:"owner_id,omitempty"`
}

type Task struct {
	gorm.Model
	Project_id  uint      `json:"project_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Due_date    time.Time `json:"due_date,omitempty"`
}

type Comment struct {
	gorm.Model
	Task_id uint   `json:"task_id,omitempty"`
	User_id uint   `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}

type Resproject struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}
