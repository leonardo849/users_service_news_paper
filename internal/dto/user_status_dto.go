package dto

import "github.com/google/uuid"

type CreateUserStatusDTO struct {
	UserId uuid.UUID `json:"user_id"`
}