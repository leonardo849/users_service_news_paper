package router

import (
	"time"
	"users-service/internal/handler"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/middleware"
	"users-service/internal/repository"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	redisLib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func desactiveCodesJob(db *gorm.DB, rc *redis.Client) {
	if db == nil {
		logger.ZapLogger.Fatal("repository.db is nill")
	}
	userRepository := repository.CreateUserRepository(db)
	userStatusRepository := repository.CreateUserStatusRepository(db)
	appRepository := repository.CreateAppRepository(userRepository, userStatusRepository, db)
	userServiceRedis := repository.CreateUserRepositoryRedis(rc)
	userService := service.CreateUserService(db, userServiceRedis, userStatusRepository, userRepository, appRepository)
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
	desactiveCodesJob(db, rc)

	userRepository := repository.CreateUserRepository(db)
	userStatusRepository := repository.CreateUserStatusRepository(db)
	appRepository := repository.CreateAppRepository(userRepository, userStatusRepository, db)
	userRepositoryRedis := repository.CreateUserRepositoryRedis(rc)
	userService := service.CreateUserService(db, userRepositoryRedis, userStatusRepository, userRepository, appRepository)
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