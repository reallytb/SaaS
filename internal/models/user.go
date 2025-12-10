package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `json:"email,omitempty" gorm:"unique"`
	PasswordHash string `json:"password_hash,omitempty"`
	Name         string `json:"name,omitempty"`
}
