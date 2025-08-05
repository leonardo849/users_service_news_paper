package router

import (
	"users-service/internal/dto"
	"os"
	"users-service/internal/logger"
	"users-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func SetupApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	
	logger.ZapLogger.Info("cors is ready")
	app.Use(middleware.LogRequestsMiddleware())

	usersGroup := app.Group("/users")
	setupUserRoutes(usersGroup)

	// @Summary Hello
	// @Description welcome message
	// @Accept json
	// @Produce json
	// @Success 200 {object} dto.MessageDTO
	// @Router / [get]
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(dto.MessageDTO{Message: "what's up"})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	logger.ZapLogger.Info("swagger is ready")

	logger.ZapLogger.Info("app is running!")
	return  app
}

func RunServer() error {
	app := SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	return app.Listen(":" + port)
}