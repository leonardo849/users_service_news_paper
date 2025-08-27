package middleware

import (
	"fmt"
	"users-service/internal/logger"
	"github.com/thoas/go-funk"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SameIdOrRole(roles []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		id := ctx.Params("id")
		idJwt := user["id"].(string)
		role := user["role"].(string)
		logger.ZapLogger.Info(fmt.Sprintf("searched id: %s. user's id: %s", id, idJwt))
		if id != idJwt && !funk.Contains(roles, role){
			return ctx.Status(403).JSON(fiber.Map{"error": "searched id and your id aren't the same id"})
		}
		return ctx.Next()
	}
}