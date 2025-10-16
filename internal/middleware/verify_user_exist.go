package middleware

import (
	"time"
	"users-service/internal/dto"
	"users-service/internal/logger"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/service"
	"users-service/internal/unitofwork"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var userRepository = repository.CreateUserRepository(repository.DB)
var userStatusRepository = repository.CreateUserStatusRepository(repository.DB)
var userRepositoryRedis = repository.CreateUserRepositoryRedis(redis.Rc)
var unitOfWork = unitofwork.CreateUnitOfWork(userRepository, userStatusRepository, repository.DB)
var userService = service.CreateUserService(nil, userRepositoryRedis, userStatusRepository, userRepository, unitOfWork)

func VerifyIfUserExistsAndIfUserIsExpired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userRepository.SetDB(repository.DB)
		userStatusRepository.SetDB(repository.DB)
		userRepositoryRedis.SetRedisDB(redis.Rc)
		unitOfWork.SetDb(repository.DB)
		userService.SetDB(repository.DB)
		

		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		id := user["id"].(string)
		logger.ZapLogger.Info("user id: " + id)
		status, reply := userService.FindOneUserById(id, ctx.Context())
		if status == 404 {
			return ctx.Status(401).JSON(fiber.Map{"error": "your user doesn't exist"})
		}
		updatedAt := user["updatedAt"]
		timestamp, ok := updatedAt.(string)
		if !ok {
			return ctx.Status(401).JSON(fiber.Map{"error": "token is wrong"})
		}

		updatedAtInTime, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{"error": "token is wrong"})
		}
		updatedAtInTime = updatedAtInTime.Truncate(time.Second)
		updatedAtInDB := reply.(dto.FindUserDTO).UpdatedAt.Truncate(time.Second)
		log.Print(updatedAtInDB, updatedAtInTime)
		if updatedAtInTime.Before(updatedAtInDB) {
			return ctx.Status(401).JSON(fiber.Map{"error": "that token is expired"})
		}
		return ctx.Next()
	}
}
