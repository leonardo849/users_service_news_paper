package router

import (
	_ "template-backend/internal/dto"
	"os"
	"template-backend/internal/logger"
	"template-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func SetupApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	
	logger.ZapLogger.Info("cors is ready")
	app.Use(middleware.LogRequestsMiddleware())
	
	// @Summary Hello
	// @Description welcome message
	// @Accept json
	// @Produce json
	// @Sucess 200 {object} dto.MessageDTO
	// @Router / [get]
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{"message": "what's up?"})
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
		port = "8080"
	}

	return app.Listen(":" + port)
}