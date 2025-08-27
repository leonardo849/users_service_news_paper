package router

import (
	"users-service/internal/handler"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/middleware"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(userGroup fiber.Router) {
	userServiceRedis := service.CreateUserServiceRedis(redis.Rc)
	userService := service.CreateUserService(repository.DB, userServiceRedis)
	userController := handler.UserController{UserService: userService}
	userGroup.Post("/create", userController.CreateUser())
	userGroup.Post("/login", userController.LoginUser())
	userGroup.Get("/one/:id", middleware.VerifyJWT(),middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.SameIdOrRole([]string{helper.Ceo, helper.Master}) ,userController.FindOneUser())
	userGroup.Put("/update/:id", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.SameIdOrRole([]string{}), userController.UpdateOneUser())
	logger.ZapLogger.Info("user routes are running!")
}