package dto

type CreateUserDTO struct {
	Username string `json:"username" validate:"required,max=50"`
	Email string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,strongpassword"`
	Fullname string `json:"fullname" validate:"required,max=100"`
}