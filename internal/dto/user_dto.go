package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserDTO struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,strongpassword"`
	Fullname string `json:"fullname" validate:"required,max=100"`
	Role *string `json:"role" validate:"omitempty,role"`
}

type FindUserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
	IsVerified bool `json:"is_verified"`
	Role      string    `json:"role"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyCodeDTO struct {
	Code string `json:"code" validate:"required,len=6,numeric"`
}

type UpdateUserDTO struct {
	Username *string `json:"username" validate:"omitempty,max=50"`
	Email    *string `json:"email" validate:"omitempty,max=100,email"`
	Fullname *string `json:"fullname" validate:"omitempty,max=100"`
}

type UpdateUserRoleDTO struct {
	Role string `json:"role" validate:"required,role"`
}

type CreateUserFromJsonFileDTO struct {
	Username string `json:"username" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,strongpassword"`
	Fullname string `json:"fullname" validate:"required,max=100"`
	Role string `json:"role" validate:"required,role"`
}

