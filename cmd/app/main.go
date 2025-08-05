package main

import (
	"log"
	"os"
	"users-service/config"
	"users-service/internal/logger"
	"users-service/internal/repository"
	"users-service/internal/router"
	"users-service/internal/validate"

	"go.uber.org/zap"
	_ "users-service/docs"
)

// @title Backend Portfolio API
// @version 1.0
// @description api for a user services for a newspaper
// @host localhost:8081
// @BasePath /
func main() {
	if err := config.SetupEnvVar(); err != nil {
		log.Fatal(err.Error())
	}
	if err := logger.StartLogger(); err != nil {
		log.Fatal(err.Error())
	}
	if _,err := repository.ConnectToDatabase(); err != nil {
		logger.ZapLogger.Error("error in repository.connectodatabase", zap.String("function", "repository.ConnectToDatabase()"), zap.Error(err))
		os.Exit(1)
	}

	validate.StartValidator()
	if err := router.RunServer(); err != nil {
		logger.ZapLogger.Error("error in run server", 
		zap.Error(err),
		zap.String("function", "router.RunServer()"),
		)
	}
}