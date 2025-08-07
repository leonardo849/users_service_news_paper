package router

import (
	"users-service/internal/handler"
	"users-service/internal/logger"
	"users-service/internal/middleware"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(userGroup fiber.Router) {
	userServiceRedis := &service.UserServiceRedis{RC: redis.Rc}
	userService := service.UserService{DB: repository.DB, UserServiceRedis: userServiceRedis}
	userController := handler.UserController{UserService: &userService}
	userGroup.Post("/create", userController.CreateUser())
	userGroup.Post("/login", userController.LoginUser())
	userGroup.Get("/one/:id", middleware.VerifyJWT(),middleware.VerifyIfUserExistsAndIfUserIsExpired() ,userController.FindOneUser())
	logger.ZapLogger.Info("user routes are running!")
}