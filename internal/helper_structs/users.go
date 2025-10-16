package helper_structs

import (
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.CreateUserRepository(db)
}

func CreateUserRedisRepository(rc *redis.Client) *repository.UserRedisRepository {
	return repository.CreateUserRepositoryRedis(rc)
}

func CreateUserStatusRepository(db *gorm.DB) *repository.UserStatusRepository {
	return  repository.CreateUserStatusRepository(db)
}

func CreateUserService(rc *redis.Client, db *gorm.DB) *service.UserService {
	return service.CreateUserService(
		db,
		CreateUserRedisRepository(rc),
		CreateUserStatusRepository(db),
		CreateUserRepository(db),
		CreateUnitOfWork(db, CreateUserStatusRepository(db), CreateUserRepository(db)),
	)
}