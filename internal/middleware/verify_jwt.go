package middleware

import (
	"fmt"
	"strings"
	"users-service/config"
	

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)



func VerifyJWT() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return  ctx.Status(401).JSON(fiber.Map{"error": "there isn't token"})
		}
		parts := strings.Split(authHeader, " ")
		if parts[0] != "Bearer" {
			return  ctx.Status(401).JSON(fiber.Map{"error": "the token is without the prefix 'bearer'"})
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return  nil, fmt.Errorf("unexpected signing method")
			} else {
				return []byte(config.Key), nil
			}
		})
		if err != nil || !token.Valid {
			return ctx.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		claims := token.Claims.(jwt.MapClaims)
		ctx.Locals("user", claims)
		return  ctx.Next()

	}
}