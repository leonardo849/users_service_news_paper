package middleware

import (
	"users-service/internal/logger"

	"github.com/gofiber/fiber/v2"
)

func LogRequestsMiddleware() fiber.Handler {
	return  func(ctx *fiber.Ctx) error {
		method := ctx.Method()
		ip := ctx.IP()
		message := "ip:" + " " + ip + " " + "method" + " " + method
		logger.ZapLogger.Info(message)
		return ctx.Next()
	}
}