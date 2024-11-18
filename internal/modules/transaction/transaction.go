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
	trQueue		TransactionQueue
	clStorage	ClientStorage
}

type TransactionSaver interface {
	AddTransaction(ctx context.Context, transaction models.Transaction) (int, error)
	UpdateTransactionStatus(ctx context.Context, transactionID int, status string) error
	ProcessTransaction(ctx context.Context, transaction models.Transaction) error
}

type TransactionProvider interface {
	GetTransaction(ctx context.Context, transactionID int) (models.Transaction, error)
	GetTransactionsByClientID(ctx context.Context, clientID int) ([]models.Transaction, error)
	GetTransactionsByStatus(ctx context.Context, status string) ([]models.Transaction, error)
	GetTransactionsByStatusAndClientID(ctx context.Context, clientID int, status string) ([]models.Transaction, error)
}

type ClientStorage interface {
	GetClientBalance(ctx context.Context, clientID int) (float64, error)
} 

type TransactionQueue interface {
	Publish(ctx context.Context, transaction models.Transaction) error
	Consume(ctx context.Context, handler func(models.Transaction) error) error
}


func New(log *slog.Logger, trSaver TransactionSaver, trProvider TransactionProvider, trQueue TransactionQueue, clStorage ClientStorage) *TransactionService {
	return &TransactionService{
		log:        log,
		trSaver:    trSaver,
		trProvider: trProvider,
		trQueue: trQueue,
		clStorage: clStorage,
	}
}

func (ts *TransactionService) AddTransaction(ctx context.Context, transaction models.Transaction) (int, error) {
	const op = "TransactionService.AddTransaction"

	log := ts.log.With(slog.String("op", op))

	log.Info("Add transaction")
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
	transaction, err := ts.trProvider.GetTransaction(context.Background(), transactionID)
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
	log.Debug("info")
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

func (ts *TransactionService) GetTransactionsByStatus(ctx context.Context, status string) ([]models.Transaction, error) {
	const op = "TransactionService.GetTransactionsByStatus"

	log := ts.log.With(slog.String("op", op))

	log.Info("Get transactions by status and clientID")

	transactions, err := ts.trProvider.GetTransactionsByStatus(ctx,  status)
	if err != nil {
		log.Error("failed to get transactions by status", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Get transactions by status successfuly")

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

func (ts *TransactionService) TransferMoney(ctx context.Context, transaction models.Transaction) error {
    const op = "TransactionService.TransferMoney"
    log := ts.log.With(slog.String("op", op))

    log.Info("Starting money transfer")

    transactionID, err := ts.trSaver.AddTransaction(ctx, transaction)
    if err != nil {
        log.Error("failed to save transaction", "err", err)
        return fmt.Errorf("%s: %w", op, err)
    }
    transaction.ID = transactionID

    if err := ts.trQueue.Publish(ctx, transaction); err != nil {
        log.Error("failed to publish transaction to queue", "err", err)
        return fmt.Errorf("%s: %w", op, err)
    }
 
    log.Info("Transaction published successfully")
    return nil
}


func (ts *TransactionService) WorkerProcessTransaction(ctx context.Context) {
	const op = "TransactionService.WorkerProcessTransaction"
	log := ts.log.With(slog.String("op", op))

	handler := func(transaction models.Transaction) error {
		log.Info("Processing transaction", slog.Int("transaction_id", transaction.ID))
	
		savedTransaction, err := ts.trProvider.GetTransaction(ctx, transaction.ID)
		if err != nil {
			log.Error("failed to get transaction from database", "err", err)
			return fmt.Errorf("%s: %w", op, err)
		}
	
		if savedTransaction.Status == "completed" || savedTransaction.Status == "failed" {
			log.Info("Transaction already processed", slog.Int("transaction_id", transaction.ID))
			return nil
		}

		err = ts.trSaver.UpdateTransactionStatus(ctx, transaction.ID, "processing")
		if err != nil {
			log.Error("failed to update transaction status to published", "err", err)
			return fmt.Errorf("%s: %w", op, err)
		}

		err = ts.trSaver.ProcessTransaction(ctx, transaction)
		if err != nil {
			log.Error("failed to transfer funds", "err", err)
			ts.trSaver.UpdateTransactionStatus(ctx, transaction.ID, "failed")
			return fmt.Errorf("%s: %w", op, err)
		}

		err = ts.trSaver.UpdateTransactionStatus(ctx, transaction.ID, "completed")
		if err != nil {
			log.Error("failed to update transaction status", "err", err)
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Info("Transaction processed successfully")
		return nil
	}

	err := ts.trQueue.Consume(ctx, handler)
	if err != nil {
		log.Error("failed to consume transactions", "err", err)
	}
}

func (ts *TransactionService) ResumePendingTransactions(ctx context.Context) error {
    const op = "TransactionService.ResumePendingTransactions"
    log := ts.log.With(slog.String("op", op))

    log.Info("Resuming pending transactions")

    pendingTransactions, err := ts.trProvider.GetTransactionsByStatus(ctx, "pending")
    if err != nil {
        log.Error("failed to get pending transactions", "err", err)
        return fmt.Errorf("%s: %w", op, err)
    }

    for _, transaction := range pendingTransactions {
        if err := ts.trQueue.Publish(ctx, transaction); err != nil {
            log.Error("failed to re-publish transaction", slog.Int("transaction_id", transaction.ID), "err", err)
            continue
        }

        err := ts.trSaver.UpdateTransactionStatus(ctx, transaction.ID, "published")
        if err != nil {
            log.Error("failed to update transaction status to published", "err", err)
            continue
        }
    }

    log.Info("Pending transactions resumed")
    return nil
}
