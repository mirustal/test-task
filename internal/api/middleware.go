package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Info("Incoming request",
			"path", c.Path(),
			"method", c.Method(),
			"ip", c.IP(),
		)
		return c.Next() 
	}
}