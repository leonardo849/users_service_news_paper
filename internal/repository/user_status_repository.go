package repository

import (
	"fmt"
	"users-service/internal/helper"
	"users-service/internal/model"

	"gorm.io/gorm"
)

type UserStatusRepository struct {
	db *gorm.DB
}


func CreateUserStatusRepository(db *gorm.DB) *UserStatusRepository {
	return  &UserStatusRepository{
		db: db,
	}
}

func (u *UserStatusRepository) CreateUserStatus(input model.UserStatusModel, tx *gorm.DB) error {
	if tx == nil {
		tx = u.db
	}

	err := tx.Create(&input).Error
	if err != nil {
		return fmt.Errorf("%s: %w", helper.INTERNALSERVER, err)
	}

	return nil
}

func (u *UserStatusRepository) SetDB(db *gorm.DB) {
	if u.db == nil {
		u.db = db
	}
}