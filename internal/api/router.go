package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"bank-service/internal/config"
	"bank-service/internal/modules/client"
	"bank-service/internal/modules/transaction"
)

func NewRouter(cfg *config.Config, log *slog.Logger, clService *client.Client, trService *transaction.TransactionService) *fiber.App {
	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	app.Use(LoggerMiddleware(log))
	app.Use(compress.New())
	app.Use(limiter.New())

	api := app.Group("/api")

	clGroup := api.Group("/client")
	clGroup.Get("/:id", func(c *fiber.Ctx) error {
		return GetClientHandler(c, clService)
	})

	trGroup := api.Group("/transactions")
	trGroup.Get("/:id", func(c *fiber.Ctx) error {
		return GetTransactionHandler(c, trService)
	})

	return app
}
