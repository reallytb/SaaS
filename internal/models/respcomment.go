package models

type RespComment struct {
	ID         uint   `json:"id"`
	TaskTitle  string `json:"task_title"`
	UserName   string `json:"user_name"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
}
