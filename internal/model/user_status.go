package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatusModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"` 
	NewsService bool `gorm:"default:false"`
}

func (u *UserStatusModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return 
}