package dto

type CreatePermissionRequest struct {
	Email string `json:"email" example:"example@example.com"`
	Role  string `json:"role" example:"viewer/editor"`
}
