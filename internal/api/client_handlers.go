package api

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"bank-service/internal/models"
	"bank-service/internal/modules/client"
)


func GetClientHandler(c *fiber.Ctx, clService *client.Client) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {  
		return c.Status(400).SendString("Invalid parameter: id must be a positive integer")
	}

	client, err := clService.GetClient(context.Background(), id)
	if err != nil {
		return c.Status(500).SendString("Error retrieving client")
	}

	return c.Status(200).JSON(client)
}

func AddClientHandler(c *fiber.Ctx, clService *client.Client) error {
	var newClient models.Client

	if err := c.BodyParser(&newClient); err != nil {
		return c.Status(400).SendString("Invalid input data")
	}

	clientID, err := clService.AddClient(c.Context(), newClient)
	if err != nil {
		return c.Status(500).SendString("Error adding client")
	}

	return c.Status(201).JSON(fiber.Map{"client_id": clientID})
}

