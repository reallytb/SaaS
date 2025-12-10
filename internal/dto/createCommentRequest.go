package dto

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000" example:"Это комментарий к задаче"`
}
