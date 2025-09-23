package repository

import (
	"context"
	"users-service/internal/model"
	

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppRepository struct {
	userRepository       *UserRepository
	userStatusRepository *UserStatusRepository
	db                   *gorm.DB
}

func CreateAppRepository(userRepository *UserRepository, userStatusRepository *UserStatusRepository, db *gorm.DB) *AppRepository {
	return  &AppRepository{
		userRepository: userRepository,
		userStatusRepository: userStatusRepository,
		db: db,
	}
}

func (a *AppRepository) CreateUserAndUserStatus(input model.UserModel, fiberCtx context.Context) (idStr *string, err error) {
	var id *uuid.UUID
	err = a.db.Transaction(func(tx *gorm.DB) error {
		
		var err error
		if id,err = a.userRepository.CreateUser(input, fiberCtx, tx); err != nil {
			return  err
		}

		if err := a.userStatusRepository.CreateUserStatus(model.UserStatusModel{UserId: *id}, tx); err != nil {
			return  err
		}
		return  nil
	})
	if err == nil {
		uuidString := id.String()
		return &uuidString, nil
	} else {
		return nil, err
	}
}


func (a *AppRepository) SetDb(db *gorm.DB) {
	if a.db == nil {
		a.db = db
	}
}