package router

import (
	"time"
	"users-service/internal/handler"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/middleware"
	"users-service/internal/redis"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func desactiveCodesJob() {
	userServiceRedis := service.CreateUserServiceRedis(redis.Rc)
	userService := service.CreateUserService(repository.DB, userServiceRedis)
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.ZapLogger.Error("panic in ExpireCodes", zap.Any("error", r))
						}
					}()
					err := userService.ExpireCodes()
					if err != nil {
						logger.ZapLogger.Error("error from userservice.findexpiratedcodes", zap.Error(err))
					} else {
						logger.ZapLogger.Info("codes were expired")
					}
				}()
			}
		}
	}()
}

func setupUserRoutes(userGroup fiber.Router) {
	desactiveCodesJob()
	userServiceRedis := service.CreateUserServiceRedis(redis.Rc)
	userService := service.CreateUserService(repository.DB, userServiceRedis)
	userController := handler.UserController{UserService: userService}
	userGroup.Post("/create", userController.CreateUser())
	userGroup.Post("/login", userController.LoginUser())
	userGroup.Get("/one/:id", middleware.VerifyJWT(),middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.SameIdOrRole([]string{helper.Ceo}) ,userController.FindOneUser())
	userGroup.Put("/update/:id", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.SameIdOrRole([]string{}), userController.UpdateOneUser())
	userGroup.Patch("/update/role/:id", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.CheckRole([]string{helper.Ceo}), userController.UpdateOneUserRole())
	
	logger.ZapLogger.Info("user routes are running!")
}