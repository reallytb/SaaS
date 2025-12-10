package dto

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email" example:"example@example.com"`
	PasswordHash string `json:"password_hash" binding:"required,min=6" example:"123456"`
	Name         string `json:"name" binding:"required,min=1" example:"vasiliy1234"`
}
