package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string    `gorm:"size:50;unique;not null" json:"username"`
	Email        string    `gorm:"size:100;unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
	FullName     string    `gorm:"size:100" json:"fullname"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool `gorm:"default:true" json:"is_active"`
	Role string `gorm:"default:'CUSTOMER';not null"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
 	return
}
