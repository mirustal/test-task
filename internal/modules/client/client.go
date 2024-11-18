package client

import (
	"context"
	"fmt"
	"log/slog"

	"bank-service/internal/models"
)

type Client struct {
	log         *slog.Logger
	cliSaver    ClientSaver
	cliProvider ClientProvider
}

type ClientSaver interface {
	AddClient(ctx context.Context, client models.Client) (int, error)
}

type ClientProvider interface {
	GetClient(ctx context.Context, clientID int) (models.Client, error)
}

func New(log *slog.Logger, cliSaver ClientSaver, cliProvider ClientProvider) *Client {
	return &Client{
		log:         log,
		cliSaver:    cliSaver,
		cliProvider: cliProvider,
	}
}

func (cl *Client) GetClient(ctx context.Context, clientID int) (models.Client, error) {
	const op = "ClientService.GetClient"

	log := cl.log.With(slog.String("op", op))

	log.Info("get client")

	var client models.Client

	client, err := cl.cliProvider.GetClient(ctx, clientID)
	if err != nil {
		log.Error("failed to get client", "err", err)
		return client, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("client get successfully")
	return client, nil
}

func (cl *Client) AddClient(ctx context.Context, client models.Client) (int, error) {
	const op = "ClientService.GetClient"

	log := cl.log.With(slog.String("op", op))

	log.Info("Add client")
	clientID, err := cl.cliSaver.AddClient(ctx, client)
	if err != nil {
		log.Error("failed to add client", "err", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Add client successfuly")

	return clientID, nil
}