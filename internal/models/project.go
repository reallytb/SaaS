package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     uint   `json:"owner_id,omitempty"`
}
