package middleware

import (
	"fmt"
	"users-service/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsVerified(expectedValue bool) fiber.Handler{
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		idJwt := user["id"].(string)
		sts, message := userService.FindOneUserById(idJwt, ctx.Context())
		if sts >= 400 {
			return  ctx.Status(sts).JSON(fiber.Map{"error": message})
		}
		userDto := message.(dto.FindUserDTO)
		if userDto.IsVerified == expectedValue {
			return  ctx.Next()
		} else {
			return  ctx.Status(401).JSON(fiber.Map{"error": fmt.Sprintf("it was supposed your user to be is_active: %t", expectedValue)})
		}
	}
}