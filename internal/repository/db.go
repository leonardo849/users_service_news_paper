package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	const maxTries = 11
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
	err = createJhonDoe(db)
	if err != nil {
		logger.ZapLogger.Error("error in create jhon doe account")
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

func createJhonDoe(db *gorm.DB) error {

	chosenEmail := os.Getenv("EMAIL_JHONDOE")
	chosenPassword := os.Getenv("PASSWORD_JHONDOE")
	if chosenEmail == "" || chosenPassword == "" {
		logger.ZapLogger.Error("the jhondoe's password or jhondoe's email is empty", zap.Error(fmt.Errorf("the jhondoe's password or jhondoe's email is empty")))
		return fmt.Errorf("the jhondoe's password or jhondoe's email is empty")
	}
	ctx := context.Background()
	const username = "Jhon"
	const name = "Jhon Doe"
	_, err := gorm.G[model.UserModel](db).Where(&model.UserModel{Username: username, FullName: name, Email: chosenEmail}).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dto := dto.CreateUserDTO{
				Username: username,
				Email:    chosenEmail,
				Password: chosenPassword,
				Fullname: name,
			}
			if err := validate.Validate.Struct(dto); err != nil {
				logger.ZapLogger.Error("error in validate struct dto", zap.Error(err))
				return err
			}
			hash, err := helper.StringToHash(dto.Password)
			if err != nil {
				logger.ZapLogger.Error("error in string to hash", zap.Error(err))
				return err
			}
			user := model.UserModel{
				Username: dto.Username,
				Email:    dto.Email,
				Password: hash,
				FullName: dto.Fullname,
				Role:     helper.Master,
				IsActive: true,
			}
			err = gorm.G[model.UserModel](db).Create(ctx, &user)
			if err != nil {
				logger.ZapLogger.Error("error in create jhon doe", zap.Error(err))
				return err
			}
			logger.ZapLogger.Info("jhon doe account was created")
			return  nil
		} else {
			logger.ZapLogger.Error("error", zap.Error(err))
			return  err
		}
	}
	return nil
}
