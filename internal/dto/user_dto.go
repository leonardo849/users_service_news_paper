package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserDTO struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,strongpassword"`
	Fullname string `json:"fullname" validate:"required,max=100"`
}

type FindUserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

type LoginUserDTO struct {
	Email     string    `json:"email"`
	Password string `json:"password"`
}