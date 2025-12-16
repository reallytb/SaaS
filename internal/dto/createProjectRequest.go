package dto

type CreateProjectRequest struct {
	Name        string `json:"name" example:"Написать API"`
	Description string `json:"description" example:"описание проекта"`
}
