package models

import "time"

type TransactionHistory struct {
	ID            int       `json:"id"`
	TransactionID int       `json:"transaction_id"`
	FromClientID  int       `json:"from_client_id"`
	ToClientID    int       `json:"to_client_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`         // "pending", "completed", "failed"
	BalanceBefore float64   `json:"balance_before"` 
	BalanceAfter  float64   `json:"balance_after"`  
	EventTime     time.Time `json:"event_time"`
}
