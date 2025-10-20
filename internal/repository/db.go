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
	"github.com/leonardo849/utils_for_backend/pkg/hash"

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

func publishSeed(db *gorm.DB, usernames []string) error {
	var users []model.UserModel

	if err := db.Where("username IN ?", usernames).Find(&users).Error; err != nil {
		logger.ZapLogger.Error("error", zap.Error(err))
		return  err
	}

	mapped := funk.Map(users, func(u model.UserModel) dtoSl.AuthPublishUserCreated {
		return dtoSl.AuthPublishUserCreated{
			AuthId: u.ID.String(),
			Username: u.Username,
			Role: u.Role,
		}
	}).([]dtoSl.AuthPublishUserCreated)

	if err := rabbitmq.GetRabbitMQClient().PublishSeed(mapped, context.Background()); err != nil {
		logger.ZapLogger.Error("error", zap.Error(err))
		return err
	}
	return  nil
}

func createAccounts(db *gorm.DB) error {
	var users []dto.CreateUserFromJsonFileDTO

	projectRoot := config.FindProjectRoot()
	if projectRoot == "" {
		return os.ErrNotExist
	}

	envPath := filepath.Join(projectRoot, "config", "users.json")
	data, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	var usersModel []*model.UserModel

	for _, newUser := range users {
		if err := validate.Validate.Struct(newUser); err != nil {
			logger.ZapLogger.Warn("invalid user data", zap.String("username", newUser.Username), zap.Error(err))
			continue
		}

		hash, err := hash.StringToHash(newUser.Password)
		if err != nil {
			logger.ZapLogger.Error("error hashing password", zap.Error(err))
			continue
		}

		newUserModel := model.UserModel{
			Username:   newUser.Username,
			Password:   hash,
			Email:      newUser.Email,
			FullName:   newUser.Fullname,
			Role:       newUser.Role,
			IsActive:   true,
			CodeDate:   nil,
			IsVerified: true,
		}


		_, err = gorm.G[model.UserModel](db).Where("username = ? OR email = ?", newUser.Username, newUser.Email).First(context.Background())

		if err == nil {
	
			logger.ZapLogger.Info("user already exists", zap.String("username", newUserModel.Username))
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error checking if user exist", zap.Error(err))
			continue
		}

		
		usersModel = append(usersModel, &newUserModel)
	}

	
	if len(usersModel) > 0 {
		if err := db.Create(usersModel).Error; err != nil {
			logger.ZapLogger.Error("error creating users", zap.Error(err))
			return err
		}
		logger.ZapLogger.Info("new users created", zap.Int("count", len(usersModel)))
	}

	
	usernames := funk.Map(users, func(u dto.CreateUserFromJsonFileDTO) string {
		return u.Username
	}).([]string)

	if err := publishSeed(db, usernames); err != nil {
		logger.ZapLogger.Error("error publishing seed", zap.Error(err))
		return err
	}

	return nil
}
