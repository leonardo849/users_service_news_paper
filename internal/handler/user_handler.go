package handler

import (
	"users-service/internal/dto"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController struct {
	UserService *service.UserService
}

func (u *UserController) CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateUserDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in create user user controller", zap.Error(err), zap.String("function", "usercontroller.createuser"))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, message := u.UserService.CreateUser(input, ctx.Context())
		property := helper.SetPropertyRequest(status)
		logger.ZapLogger.Info("returning reply create user")
		return  ctx.Status(status).JSON(fiber.Map{property: message})
	}
}