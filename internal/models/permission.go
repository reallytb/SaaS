package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Project_id uint   `json:"project_id,omitempty"`
	User_id    uint   `json:"user_id,omitempty"`
	Role       string `json:"role,omitempty"`
}
