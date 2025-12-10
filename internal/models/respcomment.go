package models

type RespComment struct {
	ID        uint   `json:"id" example:"1"`
	TaskTitle string `json:"task_title" example:"Реализовать API"`
	UserName  string `json:"user_name" example:"Иван Иванов"`
	Content   string `json:"content" example:"Это комментарий к задаче"`
	CreatedAt string `json:"created_at" example:"2023-12-15 14:30:45"`
}
