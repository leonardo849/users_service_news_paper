package redis

import (
	"fmt"
	"os"
	"strconv"
	"users-service/internal/logger"

	"context"
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
	pong, err := rc.Ping(context.Background()).Result()
	if err != nil {
		logger.ZapLogger.Fatal("error connecting to redis", zap.Error(err))
	} else {
		logger.ZapLogger.Info("redis is connected pong: " + pong )
	}
	Rc = rc
	return rc, nil
}