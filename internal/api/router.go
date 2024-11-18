package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"bank-service/internal/modules/client"
	"bank-service/internal/modules/transaction"
	"bank-service/pkg/config"
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
	clGroup.Post("/", func(c *fiber.Ctx) error {
		return AddClientHandler(c, clService)
	})

	trGroup := api.Group("/transactions")
	trGroup.Get("/:id", func(c *fiber.Ctx) error {
		return GetTransactionHandler(c, trService)
	})
	trGroup.Post("", func(c *fiber.Ctx) error {
		return AddTransactionHandler(c, trService)
	})
	trGroup.Get("/:id", func(c *fiber.Ctx) error {
		return GetClientTransactionsHandler(c, trService)
	})
	trGroup.Get("/idstatus/:id:status", func(c *fiber.Ctx) error {
		return GetTransactionsByStatusHandler(c, trService)
	})
	trGroup.Get("/status/:id", func(c *fiber.Ctx) error {
		return UpdateTransactionStatusHandler(c, trService)
	})

	return app
}
