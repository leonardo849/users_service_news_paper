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
	rc *redis.Client
}

func CreateUserServiceRedis(rc *redis.Client) *UserServiceRedis {
	return &UserServiceRedis{
		rc: rc,
	}
}

func (u *UserServiceRedis) getKey(id string) string {
	return "user" + ":" + id
}

func (u *UserServiceRedis) SetUser(user dto.FindUserDTO, fiberCtx context.Context) error {
	key := u.getKey(user.ID.String())
	logger.ZapLogger.Info(key)
	json, err := json.Marshal(user)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal(user dto.FindUserDTO)", zap.Error(err), zap.String("function", "userServiceRedis.SetUser"))
		return err
	}
	if redisStatus := u.rc.Set(fiberCtx, key, json, 30*time.Minute); redisStatus.Err() != nil {
		logger.ZapLogger.Error("error in redisClient.set", zap.Error(err), zap.String("function", "userServiceRedis.SetUser"))
		return redisStatus.Err()
	}

	logger.ZapLogger.Info("user was setted in redis")
	return nil
}

func (u *UserServiceRedis) IsRedisNil() bool {
	return u.rc == nil
}

func (u *UserServiceRedis) SetRedisDB(redisClient *redis.Client) {
	if u.IsRedisNil() {
		u.rc = redisClient
	}
}
