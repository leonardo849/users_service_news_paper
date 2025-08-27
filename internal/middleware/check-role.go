package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
)

func CheckRole(roles []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		role := user["role"].(string)
		if !funk.Contains(roles, role) {
			return ctx.Status(403).JSON(fiber.Map{"error": "your role can't do it"})
		} 
		return  ctx.Next()
	}
}