package router

import (
	"time"
	"users-service/config"
	"users-service/internal/handler"
	"users-service/internal/helper"
	"users-service/internal/logger"
	"users-service/internal/middleware"
	"users-service/internal/repository"
	"users-service/internal/service"
	"users-service/internal/unitofwork"

	"github.com/gofiber/fiber/v2"
	middlewareSl "github.com/leonardo849/shared_library_news_paper/pkg/middlewares"
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
	unitOfWork := unitofwork.CreateUnitOfWork(userRepository, userStatusRepository, db)
	userServiceRedis := repository.CreateUserRepositoryRedis(rc)
	userService := service.CreateUserService(db, userServiceRedis, userStatusRepository, userRepository, unitOfWork)
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
	unitOfWork := unitofwork.CreateUnitOfWork(userRepository, userStatusRepository, db)
	userRepositoryRedis := repository.CreateUserRepositoryRedis(rc)
	userService := service.CreateUserService(db, userRepositoryRedis, userStatusRepository, userRepository, unitOfWork)
	userController := handler.UserController{UserService: userService}
	userGroup.Get("/all", middlewareSl.VerifyJWT(config.Key), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsActiveOrInactive(true)  , middlewareSl.CheckRole([]string{helper.Ceo}) ,userController.FindAllUsers())
	userGroup.Post("/create", userController.CreateUser())
	userGroup.Get("/new_code", middlewareSl.VerifyJWT(config.Key), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsVerified(false), userController.GetNewCode())
	userGroup.Post("/verify", middlewareSl.VerifyJWT(config.Key), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsVerified(false), userController.VerifyUser())
	userGroup.Post("/login", userController.LoginUser())
	userGroup.Get("/one/:id", middlewareSl.VerifyJWT(config.Key),middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsActiveOrInactive(true) ,middlewareSl.SameIdOrRole([]string{helper.Ceo})  ,userController.FindOneUser())
	userGroup.Put("/update/:id", middlewareSl.VerifyJWT(config.Key), middleware.VerifyIfUserExistsAndIfUserIsExpired(),middleware.IsActiveOrInactive(true) , middlewareSl.SameIdOrRole([]string{}), userController.UpdateOneUser())
	userGroup.Patch("/update/role/:id", middlewareSl.VerifyJWT(config.Key), middleware.VerifyIfUserExistsAndIfUserIsExpired(), middleware.IsActiveOrInactive(true) , middlewareSl.CheckRole([]string{helper.Ceo}), userController.UpdateOneUserRole())
	
	logger.ZapLogger.Info("user routes are running!")
}