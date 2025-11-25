package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Task_id uint   `json:"task_id,omitempty"`
	User_id uint   `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}
