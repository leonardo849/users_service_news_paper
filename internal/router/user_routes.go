package router

import (
	"users-service/internal/handler"
	"users-service/internal/logger"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(userGroup fiber.Router) {
	userService := service.UserService{DB: repository.DB}
	userController := handler.UserController{UserService: &userService}
	userGroup.Post("/create", userController.CreateUser())
	logger.ZapLogger.Info("user routes are running!")
}