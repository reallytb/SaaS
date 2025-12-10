package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	TaskId  uint   `json:"task_id,omitempty"`
	UserId  uint   `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}
