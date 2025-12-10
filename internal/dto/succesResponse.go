package dto

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Операция выполнена успешно"`
	Data    interface{} `json:"data,omitempty"`
}
