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
	db               *gorm.DB
	userServiceRedis *UserServiceRedis
}

func CreateUserService(db *gorm.DB, userServiceRedis *UserServiceRedis) *UserService {
	return &UserService{
		db:               db,
		userServiceRedis: userServiceRedis,
	}
}

func (u *UserService) CreateUser(input dto.CreateUserDTO, fiberCtx context.Context) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		logger.ZapLogger.Error("validate error dto.createuserdto", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 400, err.Error()
	}
	var newUser model.UserModel
	hash, err := helper.StringToHash(input.Password)
	if err != nil {
		logger.ZapLogger.Error("internal server in stringToHash", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	_, err = gorm.G[model.UserModel](u.db).Where("username = ?", input.Username).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that username", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 409, "there is already a user with that username"
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server in find by username", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	_, err = gorm.G[model.UserModel](u.db).Where("email = ?", input.Email).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that email", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 409, "there is already a user with that email"
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server in find by email", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	newUser = model.UserModel{
		Username: input.Username,
		Email:    input.Email,
		Password: hash,
		FullName: input.Fullname,
	}

	if err = gorm.G[model.UserModel](u.db).Create(fiberCtx, &newUser); err != nil {
		logger.ZapLogger.Error("internal server in create newuser", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	msg := "user was created"
	if err = u.userServiceRedis.SetUser(dto.FindUserDTO{ID: newUser.ID, Username: newUser.Username, Email: newUser.Email, FullName: newUser.FullName, CreatedAt: newUser.CreatedAt, UpdatedAt: newUser.UpdatedAt, IsActive: newUser.IsActive, Role: newUser.Role}, fiberCtx); err != nil {
		logger.ZapLogger.Error("error in set user in database", zap.String("function", "userService.CreateUser"), zap.Error(err))
		msg = "user was created, but user wasn't setted in cache"
	}
	logger.ZapLogger.Info("new user was created")
	m := map[string]string{
		"message": msg,
		"id":      newUser.ID.String(),
	}

	return 201, m
}

func (u *UserService) FindOneUser(id string, fiberCtx context.Context) (status int, message interface{}) {
	
	userRedis, err := u.userServiceRedis.FindUser(id, fiberCtx)
	if err == nil {
		logger.ZapLogger.Info("user was gotten from redis")
		return 200, *userRedis
	}

	user, err := gorm.G[model.UserModel](u.db).Where("id = ?", id).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with id "  + id + " doesn't exist", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return 404, "user with that id doesn't exists"
		} else {
			logger.ZapLogger.Error("internal server", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return 500, err.Error()
		}
	}
	dto := dto.FindUserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsActive:  user.IsActive,
		Role: user.Role,
	}
	err = u.userServiceRedis.SetUser(dto, fiberCtx)
	msg := "returning dto"
	if err != nil {
		message = "returning dto without cache"
	}
	logger.ZapLogger.Info(msg)
	return 200, dto
}

func (u *UserService) LoginUser(dto dto.LoginUserDTO, fiberCtx context.Context) (status int, message interface{}) {
	user, err := gorm.G[model.UserModel](u.db).Where("email = ?", dto.Email).First(fiberCtx)
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

	jwt, err := helper.GenerateJWT(user.ID.String(), user.UpdatedAt, user.Email, user.Role)
	if err != nil {
		logger.ZapLogger.Error("internal server in generate jwt", zap.String("function", "userService.loginuser"), zap.Error(err))
		return 500, err.Error()
	}
	return 200, map[string]string{
		"token": jwt,
	}
}

func (u *UserService) UpdateUser(input dto.UpdateUserDTO, fiberCtx context.Context ,id string) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		logger.ZapLogger.Error("error in validate input", zap.Error(err))
		return 400, err.Error()
	}

	fields := map[string]interface{}{}

	if input.Username != nil {
		fields["username"] = *input.Username
	}
	if input.Email != nil {
		fields["email"] = *input.Email
	}
	if input.Fullname != nil {
		fields["fullname"] = *input.Fullname
	}


	if result := u.db.Model(&model.UserModel{}).Where("id = ?", id).Updates(fields); result.Error != nil {
		logger.ZapLogger.Error("error in update user", zap.Error(result.Error))
		return 500, result.Error.Error()
	}
	sts, msg := u.FindOneUser(id, fiberCtx)
	if sts >= 400 {
		return sts, msg
	}
	user, ok := msg.(dto.FindUserDTO)
	if !ok {
		logger.ZapLogger.Error("error in message to dto.finduserdto")
	} else {
		if err := u.userServiceRedis.SetUser(user, fiberCtx); err != nil {
			logger.ZapLogger.Error("error in set user in redis")
		} else {
			logger.ZapLogger.Info("updated user was setted in redis")
		}
	}
	return 200, "user was updated"
}

func (u *UserService) IsDBnil() bool {
	return u.db == nil
}

func (u *UserService) SetDB(db *gorm.DB) {
	if u.IsDBnil() {
		u.db = db
	}
}