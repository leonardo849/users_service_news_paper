package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"users-service/config"
	"users-service/internal/dto"
	"users-service/internal/logger"
	"users-service/internal/model"
	"users-service/internal/rabbitmq"
	"users-service/internal/validate"
	"users-service/pkg/hash"

	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"

	"github.com/thoas/go-funk"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URI")
	var err error
	if dsn == "" {
		logger.ZapLogger.Error("there isn't dsn")
		return nil, fmt.Errorf("there isn't dsn")
	}
	const maxTries = 20
	secondDelay := os.Getenv("SECOND_DELAY")
	var secondInt int
	var db *gorm.DB
	if secondDelay == "" {
		secondInt = 1
	} else {
		secondInt, err = strconv.Atoi(secondDelay)
		if err != nil {
			logger.ZapLogger.Error("second delay to int error")
			return nil, fmt.Errorf("second delay to int error")
		}
	}
	for i := 0; i < maxTries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.ZapLogger.Error("error in connect to db", zap.Error(err))
			if i == maxTries - 1 {
				return  nil, err
			}
			time.Sleep(time.Second * time.Duration(secondInt))
			
		} else {
			break
		}

	}

	DB = db
	err = migrateModels(db)
	if err != nil {
		logger.ZapLogger.Error("serror in migrate models")
		return nil, err
	}
	err = createAccounts(db)
	if err != nil {
		logger.ZapLogger.Error("error in create accounts")
		return  nil, err
	}
	return db, nil
}

func migrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&model.UserModel{}, &model.UserStatusModel{})
	if err != nil {
		return err
	}
	logger.ZapLogger.Info("tables are ok")
	return nil
}

func createAccounts(db *gorm.DB) error {

	var users []dto.CreateUserFromJsonFileDTO
 	projectRoot := config.FindProjectRoot()
	if projectRoot == "" {
		return  os.ErrNotExist
	}
	envPath := filepath.Join(projectRoot, "config", "users.json")
	data, err := os.ReadFile(envPath)
	if err != nil {
		return  err
	}
	if err := json.Unmarshal(data, &users); err != nil {
		return  err
	}
	var usersModel []*model.UserModel
	for i := 0; i < len(users); i++ {
		newUser := users[i]
		if err := validate.Validate.Struct(newUser); err != nil {
			return  err
		}
		hash, err := hash.StringToHash(newUser.Password)
		if err != nil {
			return  err
		}
		newUserModel := model.UserModel{
			Username: newUser.Username,
			Password: hash,
			Email: newUser.Email,
			FullName: newUser.Fullname,
			Role: newUser.Role,
			IsActive: true,
			CodeDate: nil,
			IsVerified: true,
		}
		_, err = gorm.G[model.UserModel](db).Where("username = ? OR email = ?", newUser.Username, newUser.Email).First(context.Background())
		if err == nil {
			logger.ZapLogger.Error(fmt.Sprintf("user with username %s or email %s were created", newUserModel.Username, newUserModel.Email), zap.Error(err),zap.String("function", "createAccounts"))
			continue
		} else if !errors.Is(err ,gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error in first", zap.Error(err), zap.String("function", "createaccounts"))
			return  err
		}

		usersModel = append(usersModel, &newUserModel)
	}

	if len(usersModel) > 0 {
		if err := db.Create(usersModel).Error; err != nil {
			logger.ZapLogger.Error("error in create usersmodel", zap.Error(err), zap.String("function", "createaccounts"))
			return err
		}
		logger.ZapLogger.Info("users were created")
		mapped := funk.Map(usersModel, func (u *model.UserModel) dtoSl.AuthPublishUserCreated {
			return  dtoSl.AuthPublishUserCreated{
				AuthId: u.ID.String(),
				Username: u.Username,
				Role: u.Role,
			}
		}).([]dtoSl.AuthPublishUserCreated)
		if err := rabbitmq.GetRabbitMQClient().PublishUsersVerified(mapped, context.Background()); err != nil {
			logger.ZapLogger.Fatal(err.Error())
			return err
		}
		return  nil
	}
	
	return  nil
}
