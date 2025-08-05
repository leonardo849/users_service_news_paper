package middleware

import (
	"time"
	"users-service/internal/prometheus"

	"fmt"
	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware() fiber.Handler {
    return  func(c *fiber.Ctx) error {
        start := time.Now()
        err := c.Next() 
        duration := time.Since(start)

        status := c.Response().StatusCode()
        path := c.Route().Path 
        method := c.Method()

        prometheus.HttpRequestsTotal.WithLabelValues(method, path, fmt.Sprint(status)).Inc()
        prometheus.HttpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

        return err
    }
}