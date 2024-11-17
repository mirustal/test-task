package client

import (
	"log/slog"

	"bank-service/internal/models"
)

type Client struct {
	log         *slog.Logger
	cliSaver    ClientSaver
	cliProvider ClientProvider
}

type ClientSaver interface {
}

type ClientProvider interface {
}

func New(log *slog.Logger, cliSaver ClientSaver, cliProvider ClientProvider) *Client {
	return &Client{
		log:         log,
		cliSaver:    cliSaver,
		cliProvider: cliProvider,
	}
}

func (cl *Client) GetClient(userID int) (*models.Client, error) {
	return nil, nil
}
