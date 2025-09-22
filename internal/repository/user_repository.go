package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(input model.UserModel, fiberCtx context.Context) error {

	_, err := gorm.G[model.UserModel](u.db).Where("username = ? OR email = ?", input.Username, input.Email).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that username or that email", zap.String("function", "userRepository.CreateUser"))
		return fmt.Errorf("%s: user with that username or that email", helper.CONFLICT)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server error", zap.Error(err))
		return fmt.Errorf("%s: %s", helper.INTERNALSERVER, err.Error())
	}

	
	err = gorm.G[model.UserModel](u.db).Create(fiberCtx, &input)
	if err != nil {
		logger.ZapLogger.Error("internal server error", zap.Error(err))
		return fmt.Errorf("%s: %s", helper.INTERNALSERVER, err.Error())
	}

	return nil
}

func (u *UserRepository) SetDB(db *gorm.DB) {
	if u.db == nil {
		u.db = db
	}
}

func (u *UserRepository) ExpireCodes() error {
	cutoff := time.Now().Add(-5 * time.Minute)
	result := u.db.Model(&model.UserModel{}).Where("code_date <= ? AND code IS NOT NULL", cutoff).Updates(map[string]interface{}{"code": nil, "code_date": nil})
	if result.Error != nil {
		logger.ZapLogger.Error("error in find expirated codes", zap.Error(result.Error))
		return result.Error
	}
	logger.ZapLogger.Info("codes were expired")
	return  nil
}