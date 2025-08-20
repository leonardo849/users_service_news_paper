package repository

import (
	"fmt"
	"os"
	"strconv"
	"users-service/internal/logger"
	"users-service/internal/model"
	_ "users-service/internal/model"

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
	const maxTries = 10
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
	for i := 0; i<maxTries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.ZapLogger.Error("error in connect to db", zap.Error(err))
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