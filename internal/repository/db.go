package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"users-service/config"
	"users-service/internal/dto"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/model"
	"users-service/internal/validate"

	"time"

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
	err := db.AutoMigrate(&model.UserModel{})
	if err != nil {
		return err
	}
	logger.ZapLogger.Info("tables are ok")
	return nil
}

func createAccounts(db *gorm.DB) error {

	// chosenEmail := os.Getenv("EMAIL_JHONDOE")
	// chosenPassword := os.Getenv("PASSWORD_JHONDOE")
	// if chosenEmail == "" || chosenPassword == "" {
	// 	logger.ZapLogger.Error("the jhondoe's password or jhondoe's email is empty", zap.Error(fmt.Errorf("the jhondoe's password or jhondoe's email is empty")))
	// 	return fmt.Errorf("the jhondoe's password or jhondoe's email is empty")
	// }
	// ctx := context.Background()
	// const username = "Jhon"
	// const name = "Jhon Doe"
	// _, err := gorm.G[model.UserModel](db).Where(&model.UserModel{Username: username, FullName: name, Email: chosenEmail, Role: helper.Master}).First(ctx)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		dto := dto.CreateUserDTO{
	// 			Username: username,
	// 			Email:    chosenEmail,
	// 			Password: chosenPassword,
	// 			Fullname: name,
	// 		}
	// 		if err := validate.Validate.Struct(dto); err != nil {
	// 			logger.ZapLogger.Error("error in validate struct dto", zap.Error(err))
	// 			return err
	// 		}
	// 		hash, err := helper.StringToHash(dto.Password)
	// 		if err != nil {
	// 			logger.ZapLogger.Error("error in string to hash", zap.Error(err))
	// 			return err
	// 		}
	// 		user := model.UserModel{
	// 			Username: dto.Username,
	// 			Email:    dto.Email,
	// 			Password: hash,
	// 			FullName: dto.Fullname,
	// 			Role:     helper.,
	// 			IsActive: true,
	// 		}
	// 		err = gorm.G[model.UserModel](db).Create(ctx, &user)
	// 		if err != nil {
	// 			logger.ZapLogger.Error("error in create jhon doe", zap.Error(err))
	// 			return err
	// 		}
	// 		logger.ZapLogger.Info("jhon doe account was created")
	// 		return  nil
	// 	} else {
	// 		logger.ZapLogger.Error("error", zap.Error(err))
	// 		return  err
	// 	}
	// }
	// return nil

	

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
		hash, err := helper.StringToHash(newUser.Password)
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
		return  nil
	}
	
	return  nil
}
