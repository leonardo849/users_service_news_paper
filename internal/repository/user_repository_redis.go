package repository

import (
	"context"
	"encoding/json"
	"time"
	"users-service/internal/dto"
	"users-service/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserRedisRepository struct {
	rc *redis.Client
}

func CreateUserRepositoryRedis(rc *redis.Client) *UserRedisRepository {
	return &UserRedisRepository{
		rc: rc,
	}
}

func (u *UserRedisRepository) getKey(id string) string {
	return "user" + ":" + id
}

func (u *UserRedisRepository) SetUser(user dto.FindUserDTO, fiberCtx context.Context) error {
	key := u.getKey(user.ID.String())
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

func (u *UserRedisRepository) FindUser(id string, fiberCtx context.Context) (*dto.FindUserDTO, error) {
	key := u.getKey(id)
	val, err := u.rc.Get(fiberCtx, key).Result()
	if err != nil {
		logger.ZapLogger.Error("error in get user", zap.Error(err), zap.String("function", "userServiceRedis.GetUser"))
		return nil, err
	}
	var user dto.FindUserDTO
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		logger.ZapLogger.Error("error in json unmarshal", zap.Error(err), zap.String("function", "userServiceRedis.GetUser"))
		return nil, err
	}
	logger.ZapLogger.Info("user with id" + " " + user.ID.String() + " " + "was found")
	return &user, nil
}



func (u *UserRedisRepository) SetRedisDB(redisClient *redis.Client) {
	if u.rc == nil {
		u.rc = redisClient
	}
}
