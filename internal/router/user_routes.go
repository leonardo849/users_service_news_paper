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
	redisLib "github.com/redis/go-redis/v9"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func desactiveCodesJob(db *gorm.DB) {
	if db == nil {
		logger.ZapLogger.Fatal("repository.db is nill")
	}
	userRepository := repository.CreateUserRepository(db)
	userStatusRepository := repository.CreateUserStatusRepository(db)
	userServiceRedis := repository.CreateUserServiceRedis(redis.Rc)
	userService := service.CreateUserService(db, userServiceRedis, userStatusRepository, userRepository)
	err := userService.ExpireCodes()
	if err != nil {
			logger.ZapLogger.Error("error from userservice.findexpiratedcodes", zap.Error(err))
		} else {
			logger.ZapLogger.Info("codes were expired")
	}
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

func setupUserRoutes(userGroup fiber.Router, db*gorm.DB, rc *redisLib.Client) {
	// desactiveCodesJob(db)
	if db == nil {
		logger.ZapLogger.Fatal("repository.db is nill")
	}
	userRepository := repository.CreateUserRepository(db)
	userStatusRepository := repository.CreateUserStatusRepository(db)
	userServiceRedis := repository.CreateUserServiceRedis(rc)
	userService := service.CreateUserService(db, userServiceRedis, userStatusRepository, userRepository)
	userController := handler.UserController{UserService: userService}
	userGroup.Get("/all", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsActiveOrInactive(true)  , middleware.CheckRole([]string{helper.Ceo}) ,userController.FindAllUsers())
	userGroup.Post("/create", userController.CreateUser())
	userGroup.Get("/new_code", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsVerified(false), userController.GetNewCode())
	userGroup.Post("/verify", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsVerified(false), userController.VerifyUser())
	userGroup.Post("/login", userController.LoginUser())
	userGroup.Get("/one/:id", middleware.VerifyJWT(),middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsActiveOrInactive(true) ,middleware.SameIdOrRole([]string{helper.Ceo})  ,userController.FindOneUser())
	userGroup.Put("/update/:id", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsActiveOrInactive(true) , middleware.SameIdOrRole([]string{}), userController.UpdateOneUser())
	userGroup.Patch("/update/role/:id", middleware.VerifyJWT(), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsActiveOrInactive(true) , middleware.CheckRole([]string{helper.Ceo}), userController.UpdateOneUserRole())
	
	logger.ZapLogger.Info("user routes are running!")
}