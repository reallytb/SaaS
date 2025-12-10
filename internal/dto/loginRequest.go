package dto

type LoginRequest struct {
	Email        string `json:"email" binding:"required" example:"example@example.com"`
	PasswordHash string `json:"password_hash" binding:"required,min=6" example:"123456"`
}
