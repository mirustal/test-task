package transaction

import (
	"log/slog"

	"bank-service/internal/models"
)

type TransactionService struct {
	log        *slog.Logger
	trSaver    TransactionSaver
	trProvider TransactionProvider
}

type TransactionSaver interface {
}

type TransactionProvider interface {
}

func New(log *slog.Logger, trSaver TransactionSaver, trProvider TransactionProvider) *TransactionService {
	return &TransactionService{
		log:        log,
		trSaver:    trSaver,
		trProvider: trProvider,
	}
}

func (cl *TransactionService) GetTransaction(transactionID int) (*models.Transaction, error) {
	return nil, nil
}
