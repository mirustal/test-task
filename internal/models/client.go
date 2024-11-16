package models

import (
	"time"
)

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
