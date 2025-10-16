package unitofwork

import (
	"context"
	"users-service/internal/logger"
	"users-service/internal/model"
	"users-service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Unitofwork struct {
	userRepository       *repository.UserRepository
	userStatusRepository *repository.UserStatusRepository
	db                   *gorm.DB
}

func CreateUnitOfWork(userRepository *repository.UserRepository, userStatusRepository *repository.UserStatusRepository, db *gorm.DB) *Unitofwork {
	return  &Unitofwork{
		userRepository: userRepository,
		userStatusRepository: userStatusRepository,
		db: db,
	}
}

func (u *Unitofwork) CreateUserAndUserStatus(input model.UserModel, fiberCtx context.Context) (idStr *string, err error) {
	var id *uuid.UUID
	err = u.db.Transaction(func(tx *gorm.DB) error {
		
		var err error
		if id,err = u.userRepository.CreateUser(input, fiberCtx, tx); err != nil {
			logger.ZapLogger.Error(err.Error())
			return  err
		}

		if err := u.userStatusRepository.CreateUserStatus(model.UserStatusModel{UserId: *id}, tx); err != nil {
			logger.ZapLogger.Error(err.Error())
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


func (a *Unitofwork) SetDb(db *gorm.DB) {
	if a.db == nil {
		a.db = db
	}
}