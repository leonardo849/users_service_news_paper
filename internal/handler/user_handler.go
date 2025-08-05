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

// @Summary Create new user
// @Description that method creates a new user
// @Tags user
// @Accept json
// @Produce json
// @Sucess 201 {object} dto.CreateUserDTO
// @Failure 409 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Failure 400 {object} dto.ErrorDTO
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