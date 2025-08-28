package redis

import (
	"fmt"
	"os"
	"strconv"
	"users-service/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Rc *redis.Client

func ConnectToRedis() (*redis.Client, error) {
	uriRedis := os.Getenv("REDIS_URI")
	databaseRedis := os.Getenv("REDIS_DATABASE")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if uriRedis == "" || databaseRedis == "" {
		err := fmt.Errorf("uri to redis or redis database  is empty")
		logger.ZapLogger.Error("uri to redis or redis database is empty", zap.String("function", "connectToRedis"))
		return nil, err
	}
	dbInt, err := strconv.Atoi(databaseRedis)
	if err != nil {
		logger.ZapLogger.Error("error in strconv.atoi(databaseredis)", zap.String("function", "connectToRedis"))
		return  nil, err
	}
	rc := redis.NewClient(&redis.Options{
		Addr: uriRedis,
		Password: redisPassword,
		DB: dbInt,
	})
	Rc = rc
	logger.ZapLogger.Info("connected to redis")
	return rc, nil
}