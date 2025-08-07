package service

import (
	"context"
	"encoding/json"
	"time"
	"users-service/internal/dto"
	"users-service/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserServiceRedis struct {
	RC *redis.Client
}

func (u *UserServiceRedis) SetUser(property string, user dto.FindUserDTO, fiberCtx context.Context) error {
	key := property + ":" + user.ID.String()
	logger.ZapLogger.Info(key)
	json, err := json.Marshal(user)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal(user model.usermodel)", zap.Error(err), zap.String("function", "userServiceRedis.SetUser"))
		return err
	}
	if redisStatus := u.RC.Set(fiberCtx, key, json, 30*time.Minute); redisStatus.Err() != nil {
		logger.ZapLogger.Error("error in redisClient.set", zap.Error(err), zap.String("function", "userServiceRedis.SetUser"))
		return  redisStatus.Err()
	}
	logger.ZapLogger.Info("user was setted in redis")
	return nil
}