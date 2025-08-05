package service

import (
	"context"
	"errors"
	"users-service/internal/dto"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/model"
	"users-service/internal/validate"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (u *UserService) CreateUser(dto dto.CreateUserDTO, fiberCtx context.Context) (status int, message string) {
	if err := validate.Validate.Struct(dto); err != nil {
		logger.ZapLogger.Error("validate error dto.createuserdto", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 400, err.Error()
	}
	var newUser model.UserModel
	hash, err := helper.StringToHash(dto.Password)
	if err != nil {
		logger.ZapLogger.Error("internal server in stringToHash", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	_, err = gorm.G[model.UserModel](u.DB).Where("username = ?", dto.Username).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that username", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 409, "there is already a user with that username"
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server in find by username", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	newUser = model.UserModel{
		Username: dto.Username,
		Email: dto.Email,
		Password: hash,
		FullName: dto.Fullname,
	}
	
	err = gorm.G[model.UserModel](u.DB).Create(fiberCtx, &newUser)
	if err != nil {
		logger.ZapLogger.Error("internal server in create newuser", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	logger.ZapLogger.Info("new user was created" )
	return 201, "user was created"
}