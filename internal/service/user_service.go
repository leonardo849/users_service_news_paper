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

func (u *UserService) CreateUser(dto dto.CreateUserDTO, fiberCtx context.Context) (status int, message interface{}) {
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
	_, err = gorm.G[model.UserModel](u.DB).Where("email = ?", dto.Email).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that email", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 409, "there is already a user with that email"
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server in find by email", zap.String("function", "userService.CreateUser"), zap.Error(err))
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
	m := map[string]string{
		"message": "user was created",
		"id": newUser.ID.String(),
	}

	return 201, m
}

func (u *UserService) FindOneUser(id string, fiberCtx context.Context) (status int, message interface{}) {
	user, err := gorm.G[model.UserModel](u.DB).Where("id = ?", id).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with that id doesn't exist", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return 404, "user with that id already exists"
		} else {
			logger.ZapLogger.Error("internal server", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return 500, err.Error()
		}
	}
	dto := dto.FindUserDTO{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsActive: user.IsActive,
	}
	logger.ZapLogger.Info("returning dto")
	return 200, dto
}

func (u *UserService) LoginUser(dto dto.LoginUserDTO, fiberCtx context.Context) (status int, message interface{}){
	user, err := gorm.G[model.UserModel](u.DB).Where("email = ?", dto.Email).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with that email doesn't exist", zap.Error(err), zap.String("function", "userService.loginuser"))
			return 404, "a user with that email doesn't exist"
		} else {
			logger.ZapLogger.Error("internal server", zap.Error(err), zap.String("function", "userService.loginuser"))
			return 500, err.Error()
		}
	}


	if !helper.CompareHash(dto.Password, user.Password) {
		logger.ZapLogger.Error("password is wrong", zap.String("function", "userService.loginuser"))
		return 401, "password is wrong"
	}

	jwt, err := helper.GenerateJWT(user.ID.String(),user.UpdatedAt, user.Email)
	if err != nil {
		logger.ZapLogger.Error("internal server in generate jwt", zap.String("function", "userService.loginuser"), zap.Error(err))
		return 500, err.Error()
	}
	return 200, map[string]string{
		"token": jwt,
	}
}