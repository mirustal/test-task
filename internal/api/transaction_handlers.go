package api

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"bank-service/internal/models"
	"bank-service/internal/modules/transaction"
)

func GetTransactionHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {  
		return c.Status(400).SendString("Invalid parameter: id must be a positive integer")
	}
	
	transaction, err := trService.GetTransaction(context.Background(), id)
	if err != nil {
		return c.Status(500).SendString("Error retrieving transaction")
	}
	return c.Status(200).JSON(transaction)
}

func AddTransactionHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	var tr models.Transaction
	if err := c.BodyParser(&tr); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	if tr.FromClientID <= 0 || tr.ToClientID <= 0 || tr.Amount <= 0 {
		return c.Status(400).SendString("Invalid transaction details")
	}

	transactionID, err := trService.AddTransaction(c.Context(), tr)
	if err != nil {
		return c.Status(500).SendString("Error adding transaction")
	}

	return c.Status(201).JSON(fiber.Map{"transaction_id": transactionID})
}

func GetClientTransactionsHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	clientIDStr := c.Params("id")

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil || clientID <= 0 {
		return c.Status(400).SendString("Invalid parameter: client_id must be a positive integer")
	}

	transactions, err := trService.GetTransactionsByClientID(c.Context(), clientID)
	if err != nil {
		return c.Status(500).SendString("Error retrieving transactions")
	}

	return c.Status(200).JSON(transactions)
}

func GetTransactionsByStatusAndIDHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	clientIDStr := c.Params("id")
	status := c.Query("status", "completed")


	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil || clientID <= 0 {
		return c.Status(400).SendString("Invalid parameter: client_id must be a positive integer")
	}

	if status == "" {
		return c.Status(400).SendString("Invalid parameter: status is required")
	}

	transactions, err := trService.GetTransactionsByStatusAndClientID(c.Context(), clientID, status)
	if err != nil {
		return c.Status(500).SendString("Error retrieving transactions")
	}

	return c.Status(200).JSON(transactions)
}

func GetTransactionsByStatusHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	status := c.Query("status", "completed")
	

	if status == "" {
		return c.Status(400).SendString("Invalid parameter: status is required")
	}

	transactions, err := trService.GetTransactionsByStatus(c.Context(), status)
	if err != nil {
		return c.Status(500).SendString("Error retrieving transactions")
	}

	return c.Status(200).JSON(transactions)
}

func UpdateTransactionStatusHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	idStr := c.Params("id")
	status := c.Query("status", "completed")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(400).SendString("Invalid parameter: id must be a positive integer")
	}

	if status == "" {
		return c.Status(400).SendString("Invalid parameter: status is required")
	}

	err = trService.UpdateTransactionStatus(c.Context(), id, status)
	if err != nil {
		return c.Status(500).SendString("Error updating transaction status")
	}

	return c.SendStatus(200)
}

func TransferMoneyHandler(c *fiber.Ctx, trService *transaction.TransactionService) error {
	var transaction models.Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	if transaction.FromClientID <= 0 || transaction.ToClientID <= 0 || transaction.Amount <= 0 {
		return c.Status(400).SendString("Invalid transaction details")
	}

	if transaction.FromClientID == transaction.ToClientID {
		return c.Status(400).SendString("Cannot transfer to the same client")
	}

	if err := trService.TransferMoney(context.Background(), transaction); err != nil {
		return c.Status(500).SendString("Error processing transaction")
	}

	return c.Status(200).SendString("Transaction submitted successfully")
}



