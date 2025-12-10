package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	ProjectId uint   `json:"project_id,omitempty"`
	UserId    uint   `json:"user_id,omitempty"`
	Role      string `json:"role,omitempty"`
}
