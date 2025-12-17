package dto

import "time"

type CreateTaskRequest struct {
	Title       string    `json:"title,omitempty" binding:"required,min=1" example:"написать функцию" `
	Description string    `json:"description,omitempty" example:"описание задачи"`
	Status      string    `json:"status,omitempty" example:"todo/in-progress/done"`
	Priority    string    `json:"priority,omitempty" example:"low/medium/high"`
	Due_date    time.Time `json:"due_date,omitempty" example:"2025-11-20T13:37:27+00:00"`
}
