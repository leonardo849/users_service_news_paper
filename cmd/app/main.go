package main

import (
	"log"
	"os"
	"users-service/config"
	"users-service/internal/logger"
	"users-service/internal/prometheus"
	"users-service/internal/rabbitmq"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/router"
	"users-service/internal/validate"

	_ "users-service/docs"

	"go.uber.org/zap"
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
	if _, err := redis.ConnectToRedis(); err != nil {
		logger.ZapLogger.Error("error in connect to redis", zap.String("function", "redis.ConnectToRedis"), zap.Error(err))
		os.Exit(1)
	}
	if  err := rabbitmq.ConnectToRabbitMQ(); err != nil {
		logger.ZapLogger.Error("error in connect to rabbit", zap.String("function", "rabbitmq.connectorabbitmq"), zap.Error(err))
	}

	prometheus.StartPrometheus()
	validate.StartValidator()



	if err := router.RunServer(); err != nil {
		logger.ZapLogger.Error("error in run server", 
		zap.Error(err),
		zap.String("function", "router.RunServer()"),
		)
	}
}