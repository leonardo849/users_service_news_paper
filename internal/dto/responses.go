package dto

type MessageDTO struct {
	Message string `json:"string"`
}

type CreateDTO struct {
	Message string `json:"string"`
	ID string `json:"id"`
}

type ErrorDTO struct {
	Error string `json:"error"`
}

type LoginDTO struct {
	Token string `json:"token"`
}