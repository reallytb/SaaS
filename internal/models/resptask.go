package models

type Resptask struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Due_date    string `json:"due_date,omitempty"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}
