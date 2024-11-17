package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"bank-service/internal/modules/client"
)


func GetClientHandler(c *fiber.Ctx, clService *client.Client) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {  
		return c.Status(400).SendString("Invalid parameter: id must be a positive integer")
	}

	client, err := clService.GetClient(id)
	if err != nil {
		return c.Status(500).SendString("Error retrieving client")
	}

	return c.Status(200).JSON(client)
}
