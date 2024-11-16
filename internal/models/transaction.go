package models

import "time"

type Transaction struct {
	ID           int       `json:"id"`
	FromClientID int       `json:"from_client_id"`
	ToClientID   int       `json:"to_client_id"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"` // "pending", "completed", "failed"
	CreatedAt    time.Time `json:"created_at"`
	ProcessedAt  time.Time `json:"processed_at,omitempty"` 
}
