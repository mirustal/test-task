package transaction

import (
	"context"
	"fmt"
	"log/slog"

	"bank-service/internal/models"
)

type TransactionService struct {
	log        *slog.Logger
	trSaver    TransactionSaver
	trProvider TransactionProvider
}

type TransactionSaver interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (int, error)
	UpdateTransactionStatus(ctx context.Context, transactionID int, status string) error
}

type TransactionProvider interface {
	GetTransaction(ctx context.Context, transactionID int) (models.Transaction, error)
	GetTransactionsByClientID(ctx context.Context, clientID int) ([]models.Transaction, error)
	GetTransactionsByStatusAndClientID(ctx context.Context, clientID int, status string) ([]models.Transaction, error)
}


func New(log *slog.Logger, trSaver TransactionSaver, trProvider TransactionProvider) *TransactionService {
	return &TransactionService{
		log:        log,
		trSaver:    trSaver,
		trProvider: trProvider,
	}
}

func (ts *TransactionService) AddTransaction(ctx context.Context, transaction models.Transaction) (int, error) {
	const op = "TransactionService.AddTransaction"

	log := ts.log.With(slog.String("op", op))

	log.Info("Add transaction")
	//Добавь проверку на деньги
	transactionID, err := ts.trSaver.AddTransaction(ctx, transaction)
	if err != nil {
		log.Error("failed to add transaction", "err", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Add transaction")
	return transactionID, nil
}

func (ts *TransactionService) GetTransaction(ctx context.Context, transactionID int) (models.Transaction, error) {
	const op = "TransactionSetrvice.GetTransaction"

	log := ts.log.With(slog.String("op", op))

	var transaction models.Transaction

	log.Info("Get transaction")
	transaction, err := ts.trProvider.GetTransaction(ctx, transactionID)
	if err != nil {
		log.Error("failed to get transaction", "err", err)
		return transaction, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Get transaction successfuly")

	return transaction, nil
}

func (ts *TransactionService) GetTransactionsByClientID(ctx context.Context, clientID int) ([]models.Transaction, error) {
	const op = "TransactionService.GetTransactionsByClientID"

	log := ts.log.With(slog.String("op", op))

	log.Info("Get transactions by clientID")

	transactions, err := ts.trProvider.GetTransactionsByClientID(ctx, clientID)
	if err != nil {
		log.Error("failed to get transactions for client", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Get transactions by clientID successfuly")

	return transactions, nil
}

func (ts *TransactionService) GetTransactionsByStatusAndClientID(ctx context.Context, clientID int, status string) ([]models.Transaction, error) {
	const op = "TransactionService.GetTransactionsByStatusAndClientID"

	log := ts.log.With(slog.String("op", op))

	log.Info("Get transactions by status and clientID")

	transactions, err := ts.trProvider.GetTransactionsByStatusAndClientID(ctx, clientID, status)
	if err != nil {
		log.Error("failed to get transactions by status for client", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Get transactions by status and clientID successfuly")

	return transactions, nil
}

func (ts *TransactionService) UpdateTransactionStatus(ctx context.Context, transactionID int, status string) error {
	const op = "TransactionService.UpdateTransactionStatus"

	log := ts.log.With(slog.String("op", op))

	log.Info("Update Transaction Status")

	err := ts.trSaver.UpdateTransactionStatus(ctx, transactionID, status)
	if err != nil {
		log.Error("failed to update transaction status", "err", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Update transaction status successfuly")

	return nil
}
