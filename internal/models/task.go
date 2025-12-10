package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ProjectId   uint      `json:"project_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Due_date    time.Time `json:"due_date,omitempty"`
}
