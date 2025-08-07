package middleware

import (
	"time"
	"users-service/internal/dto"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var userServiceRedis = service.CreateUserServiceRedis(nil)
var userService = service.CreateUserService(nil, userServiceRedis)
func VerifyIfUserExistsAndIfUserIsExpired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if userServiceRedis.IsRedisNil() {
			userServiceRedis.SetRedisDB(redis.Rc)
		}
		if userService.IsDBnil() {
			userService.SetDB(repository.DB)
		}

		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		status, reply := userService.FindOneUser(user["id"].(string), ctx.Context())
		if status >= 400 {
			return  ctx.Status(status).JSON(fiber.Map{"error": reply})
		}
		updatedAt := user["updatedAt"]
		timestamp, ok := updatedAt.(string)
		if !ok {
			return  ctx.Status(401).JSON(fiber.Map{"error": "token is wrong"})
		}

		updatedAtInTime, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{"error": "token is wrong"})
		}

		updatedAtInDB := reply.(dto.FindUserDTO).UpdatedAt
		
		if updatedAtInTime.Before(updatedAtInDB) {
			return  ctx.Status(401).JSON(fiber.Map{"error": "that token is expired"})
		}
		return  ctx.Next()
	}
}