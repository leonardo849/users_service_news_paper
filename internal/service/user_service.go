package service

import (
	"context"
	"errors"
	"time"
	"users-service/internal/dto"
	"users-service/internal/email"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/model"
	"users-service/internal/validate"
	"users-service/pkg/hash"
	"users-service/pkg/random"
	"users-service/pkg/date"
	"github.com/thoas/go-funk"
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
	hashPassword, err := hash.StringToHash(input.Password)
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
	code := random.EncodeToString(6)
	hashCode, err := hash.StringToHash(code)
	if err != nil {
		return 500, err.Error()
	}
	newUser = model.UserModel{
		Username: input.Username,
		Email:    input.Email,
		Password: hashPassword,
		FullName: input.Fullname,
		Code: &hashCode,
		CodeDate: date.PtrTime(time.Now()),
	}

	if err = gorm.G[model.UserModel](u.db).Create(fiberCtx, &newUser); err != nil {
		logger.ZapLogger.Error("internal server in create newuser", zap.String("function", "userService.CreateUser"), zap.Error(err))
		return 500, err.Error()
	}
	email.SendEmail(dto.SendEmailDTO{
		To: newUser.Email,
		Subject: "code auth",
		Text: code,
	})
	msg := "user was created"
	// if err = u.userServiceRedis.SetUser(dto.FindUserDTO{ID: newUser.ID, Username: newUser.Username, Email: newUser.Email, FullName: newUser.FullName, CreatedAt: newUser.CreatedAt, UpdatedAt: newUser.UpdatedAt, IsActive: newUser.IsActive, Role: newUser.Role}, fiberCtx); err != nil {
	// 	logger.ZapLogger.Error("error in set user in database", zap.String("function", "userService.CreateUser"), zap.Error(err))
	// 	msg = "user was created, but user wasn't setted in cache"
	// }
	logger.ZapLogger.Info("new user was created")
	m := map[string]string{
		"message": msg,
		"id":      newUser.ID.String(),
	}

	return 201, m
}

func (u *UserService) ExpireCodes() error {
	cutoff := time.Now().Add(-5 * time.Minute)
	result := u.db.Model(&model.UserModel{}).Where("code_date <= ? AND code IS NOT NULL", cutoff).Update("code", nil)
	if result.Error != nil {
		logger.ZapLogger.Error("error in find expirated codes", zap.Error(result.Error))
		return result.Error
	}
	// mapped := funk.Map(users, func(user model.UserModel) dto.FindUserDTO{
	// 	return dto.FindUserDTO{
	// 		ID: user.ID,
	// 		Username: user.Username,
	// 		Email: user.Email,
	// 		FullName: user.FullName,
	// 		CreatedAt: user.CreatedAt,
	// 		UpdatedAt: user.UpdatedAt,
	// 		IsActive: user.IsActive,
	// 		Role: user.Role,
	// 	}
	// }).([]dto.FindUserDTO)
	logger.ZapLogger.Info("codes were expired")
	return  nil
}

func (u *UserService) CreateNewCode(id string, fiberCtx context.Context) (status int, message interface{}) {
	code := random.EncodeToString(6)
	hashCode, err := hash.StringToHash(code)
	if err != nil {
		return 500, err.Error()
	}
	result := u.db.Model(&model.UserModel{}).Where("id = ? AND is_active = ?", id, false).Updates(model.UserModel{Code: &hashCode, CodeDate: date.PtrTime(time.Now())})
	if result.Error != nil {
		return 500, result.Error.Error()
	}
	user, err := gorm.G[model.UserModel](u.db).Where("id = ?",  id).First(fiberCtx)
	if err != nil {
		return 500, result.Error.Error()
	}
	email.SendEmail(dto.SendEmailDTO{
		To: user.Email,
		Subject: "new code",
		Text: code,
	})
	return 200, "new code was generated. It was sent to your email"
}

func (u *UserService) VerifyCode(id string, fiberCtx context.Context, input dto.VerifyCodeDTO) (status int, message interface{}) {
	user, err := gorm.G[model.UserModel](u.db).Where("id = ? AND is_active = ? AND code_date >= ?", id, false, time.Now().Add(-5 * time.Minute)).First(fiberCtx)
	if err != nil {
		return 500, err.Error()
	}
	if user.Code == nil {
		status, message = u.CreateNewCode(id, fiberCtx)
		return status, message
	}
	if err := validate.Validate.Struct(input); err != nil {
		return 400, err.Error()
	}
	if hash.CompareHash(input.Code, *user.Code) {
		result := u.db.Model(&model.UserModel{}).Where("id = ? AND is_active = ?", id, false).Updates(map[string]interface{}{"is_active": true, "code": nil, "code_date": nil})
		if result.Error != nil {
			return 500, result.Error.Error()
		} else {
			return 200, "your user is active"
		}
	} 
	return 400, "the code is wrong"
}

func (u *UserService) FindOneUserById(id string, fiberCtx context.Context) (status int, message interface{}) {

	userRedis, err := u.userServiceRedis.FindUser(id, fiberCtx)
	if err == nil {
		logger.ZapLogger.Info("user was gotten from redis")
		return 200, *userRedis
	}

	user, err := gorm.G[model.UserModel](u.db).Where("id = ?", id).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with id "+id+" doesn't exist", zap.Error(err), zap.String("function", "userservice.findoneuser"))
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
		Role:      user.Role,
	}
	err = u.userServiceRedis.SetUser(dto, fiberCtx)
	msg := "returning dto"
	if err != nil {
		message = "returning dto without cache"
	}
	logger.ZapLogger.Info(msg)
	return 200, dto
}

func (u *UserService) FindOneUserByEmail(email string, fiberCtx context.Context) (status int, message interface{}) {
	user, err := gorm.G[model.UserModel](u.db).Where("email = ?", email).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with email "  + email + " doesn't exist", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return 404, "user with that email doesn't exists"
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
	_, err = u.userServiceRedis.FindUser(dto.ID.String(), fiberCtx)
	if err != nil {
		err = u.userServiceRedis.SetUser(dto, fiberCtx)
		if err != nil {
			logger.ZapLogger.Info("returning dto without cache")
			return 200, dto
		} else {
			logger.ZapLogger.Info("user was setted in redis")
		}
	}
	logger.ZapLogger.Info("returning dto")
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

	if !hash.CompareHash(dto.Password, user.Password) {
		logger.ZapLogger.Error("password is wrong", zap.String("function", "userService.loginuser"))
		return 401, "password is wrong"
	}

	jwt, err := helper.GenerateJWT(user.ID.String(), user.UpdatedAt, user.Role)
	if err != nil {
		logger.ZapLogger.Error("internal server in generate jwt", zap.String("function", "userService.loginuser"), zap.Error(err))
		return 500, err.Error()
	}
	return 200, map[string]string{
		"token": jwt,
	}
}

func (u *UserService) FindAllUsers() (status int, users []dto.FindUserDTO) {
	var arr []model.UserModel
	result := u.db.Find(&arr)
	
	if result.Error != nil {
		return 500, nil
	}
	users = funk.Map(arr, func(user model.UserModel) dto.FindUserDTO {
		return dto.FindUserDTO{
			ID: user.ID,
			Username: user.Username,
			Email: user.Email,
			FullName: user.FullName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsActive: user.IsActive,
			Role: user.Role,
		}
	}).([]dto.FindUserDTO)
	return 200, users
}

func (u *UserService) UpdateUser(input dto.UpdateUserDTO, fiberCtx context.Context, id string) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		logger.ZapLogger.Error("error in validate input", zap.Error(err))
		return 400, err.Error()
	}

	sts, msg := u.FindOneUserById(id, fiberCtx)
	if sts >= 400 {
		return sts, msg
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

func (u *UserService) UpdateUserRole(input dto.UpdateUserRoleDTO, fiberCtx context.Context, id string) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		logger.ZapLogger.Error("error in validate struct", zap.Error(err), zap.String("function", "userservice.updateuserrole"))
		return 400, err.Error()
	}
	sts, msg := u.FindOneUserById(id, fiberCtx)
	if sts >= 400 {
		return sts, msg
	}
	err := u.db.Model(&model.UserModel{}).Where("id = ?", id).Update("role", input.Role).Error
	if err != nil {
		logger.ZapLogger.Error("error in update role", zap.Error(err), zap.String("function", "userservice.updateuserrole"))
		return 500, err.Error()
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
