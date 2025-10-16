package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/model"
	"github.com/leonardo849/utils_for_backend/pkg/date"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	errorsSl "github.com/leonardo849/utils_for_backend/pkg/errors"
)

type UserRepository struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(input model.UserModel, fiberCtx context.Context, tx *gorm.DB) (*uuid.UUID,  error) {
	helper.SetTx(&tx, u.db)
	_, err := gorm.G[model.UserModel](tx).Where("username = ? OR email = ?", input.Username, input.Email).First(fiberCtx)
	if err == nil {
		logger.ZapLogger.Error("there is already a user with that username or that email", zap.String("function", "userRepository.CreateUser"))
		return nil,fmt.Errorf("%s: user with that username or that email", helper.CONFLICT)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("internal server error", zap.Error(err))
		return nil,fmt.Errorf("%s: %s", helper.INTERNALSERVER, err.Error())
	}

	
	err = gorm.G[model.UserModel](tx).Create(fiberCtx, &input)
	if err != nil {
		logger.ZapLogger.Error("internal server error", zap.Error(err))
		return nil,fmt.Errorf("%s: %s", helper.INTERNALSERVER, err.Error())
	}
	idUuid := input.ID
	return &idUuid, nil
}

func (u *UserRepository) FindExpiredUser (id string, fiberCtx context.Context) (*model.UserModel, error) {
	user, err := gorm.G[model.UserModel](u.db).Where("id = ? AND is_verified = ? AND code_date >= ?", id, false, time.Now().Add(-5*time.Minute)).First(fiberCtx)
	if err != nil {
		return nil, nil
	}
	return  &user, nil
}

func (u *UserRepository) VerifyCode(id string) error {
	err := u.db.Model(&model.UserModel{}).Where("id = ? AND is_verified = ?", id, false).Updates(map[string]interface{}{"is_verified": true, "code": nil, "code_date": nil}).Error
	if err != nil {
		return  fmt.Errorf("[%s] %s", helper.INTERNALSERVER, err.Error())
	}
	return  nil
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

func (u *UserRepository) CreateNewCode(id string, fiberCtx context.Context, hashCode *string) (email *string, err error) {
	result := u.db.Model(&model.UserModel{}).Where("id = ? AND is_verified = ?", id, false).Updates(model.UserModel{Code: hashCode, CodeDate: date.PtrTime(time.Now())})
	if result.Error != nil {
		return nil,fmt.Errorf("%s:%s", helper.INTERNALSERVER, result.Error.Error())
	}
	user, err := gorm.G[model.UserModel](u.db).Where("id = ?", id).First(fiberCtx)
	if err != nil {
		return nil,fmt.Errorf("%s:%s", helper.INTERNALSERVER, err.Error())
	}
	return &user.Email, nil
}

func (u *UserRepository) FindUserById(id string, fiberCtx context.Context) (*model.UserModel, error) {
	user, err := gorm.G[model.UserModel](u.db).Where("id = ?", id).First(fiberCtx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("a user with id "+id+" doesn't exist", zap.Error(err), zap.String("function", "userservice.findoneuser"))
			return nil, fmt.Errorf("[%s] %s", errorsSl.NOTFOUND, "user with that id doesn't exists")
		} else {
			logger.ZapLogger.Error(err.Error(), zap.Error(err))
			return nil, fmt.Errorf("[%s] %s", errorsSl.INTERNALSERVER, err.Error())
		}
	}
	return &user, nil
}