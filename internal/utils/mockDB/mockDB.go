package mockdb

import (
	"context"
	"log"
	"time"

	"bank-service/internal/adapter/db/postgres"
	"bank-service/internal/models"
)

func SeedDataBase(ctx context.Context, storage *postgres.Storage) error {
	clients := []models.Client{
		{Name: "mirustal", Balance: 1000.00, CreatedAt: time.Now()},
		{Name: "mirustal2", Balance: 500.00, CreatedAt: time.Now()},
		{Name: "mirustal3", Balance: 750.00, CreatedAt: time.Now()},
	}
	var err error
	var clientIDs []int

	for _, client := range clients {
		clientID, err := storage.AddClient(ctx, client)
		if err != nil {
			log.Printf("Failed to seed client %s: %v", client.Name, err)
			continue
		}
		clientIDs = append(clientIDs, clientID)
		log.Printf("Added client: ID=%d, Name=%s, Balance=%.2f", clientID, client.Name, client.Balance)
	}

	if len(clientIDs) < 2 {
		log.Println("Not enough clients to create transactions. Seeding stopped.")
		return err
	}

	transactions := []models.Transaction{
		{
			FromClientID: clientIDs[0], ToClientID: clientIDs[1], Amount: 100.00,
			Status:    "completed",
			CreatedAt: time.Now(),
			ProcessedAt: time.Now(),
		},
		{
			FromClientID: clientIDs[1], ToClientID: clientIDs[0], Amount: 50.00,
			Status:    "completed",
			CreatedAt: time.Now(),
			ProcessedAt: time.Now(),
		},
		{
			FromClientID: clientIDs[0], ToClientID: clientIDs[2], Amount: 200.00,
			Status:    "pending",
			CreatedAt: time.Now(),
			ProcessedAt: time.Now(),
		},
	}

	for _, tr := range transactions {
		trID, err := storage.AddTransaction(ctx, tr)
		if err != nil {
			log.Printf("Failed to seed transaction: %+v, error: %v", tr, err)
			return err
		}
		log.Printf("Added transaction: ID=%d, FromClientID=%d, ToClientID=%d, Amount=%.2f, Status=%s",
			trID, tr.FromClientID, tr.ToClientID, tr.Amount, tr.Status)
	}

	return nil
}

