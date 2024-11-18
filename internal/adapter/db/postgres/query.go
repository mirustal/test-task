package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"bank-service/internal/models"
)

func (pg *Storage) AddClient(ctx context.Context, client models.Client) (int, error) {
	const op = "postgres.AddClient"

	query := `
        INSERT INTO clients (name, balance, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id;
    `
	var clientID int
	err := pg.db.QueryRow(ctx, query, client.Name, client.Balance).Scan(&clientID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return clientID, nil
}

func (pg *Storage) GetClient(ctx context.Context, clientID int) (models.Client, error) {
	const op = "postgres.GetClient"

	query := `
        SELECT id, name, balance, created_at
        FROM clients
        WHERE id = $1
        LIMIT 1;
    `
	var client models.Client

	err := pg.db.QueryRow(ctx, query, clientID).Scan(
		&client.ID,
		&client.Name,
		&client.Balance,
		&client.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return client, fmt.Errorf("client with id %d not found", clientID)
	} else if err != nil {
		return client, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

func (pg *Storage) GetClientBalance(ctx context.Context, clientID int) (float64, error) {
	const op = "postgres.GetClientBalance"

	query := `
        SELECT balance
        FROM clients
        WHERE id = $1;
    `
	var balance float64
	err := pg.db.QueryRow(ctx, query, clientID).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("%s: client not found: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return balance, nil
}

func (pg *Storage) ProcessTransaction(ctx context.Context, transaction models.Transaction) error {
	const op = "postgres.ProcessTransaction"

	log := pg.log.With(slog.String("op", op))

	tx, err := pg.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var fromBalance float64
	err = tx.QueryRow(ctx, `
		SELECT balance
		FROM clients
		WHERE id = $1
		FOR UPDATE;
	`, transaction.FromClientID).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if fromBalance < transaction.Amount {
		_, err = tx.Exec(ctx, `
			UPDATE transactions
			SET status = 'failed'
			WHERE id = $1;
		`, transaction.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE clients
		SET balance = balance - $1
		WHERE id = $2;
	`, transaction.Amount, transaction.FromClientID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE clients
		SET balance = balance + $1
		WHERE id = $2;
	`, transaction.Amount, transaction.ToClientID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE transactions
		SET status = 'completed', processed_at = NOW()
		WHERE id = $1;
	`, transaction.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error("Transaction panicked", "err", p)
			tx.Rollback(ctx)
		} else if err != nil {
			log.Error("Transaction failed", "err", err)
			tx.Rollback(ctx)
		} else {
			log.Info("Transaction committed successfully")
			err = tx.Commit(ctx)
		}
	}()

	return nil
}

func (pg *Storage) AddTransaction(ctx context.Context, transaction models.Transaction) (int, error) {
	const op = "postgres.AddTransaction"

	transaction.Status = "pending"
	query := `
        INSERT INTO transactions (from_client_id, to_client_id, amount, status, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id;
    `
	var transactionID int
	err := pg.db.QueryRow(ctx, query, transaction.FromClientID, transaction.ToClientID, transaction.Amount, transaction.Status).Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return transactionID, nil
}

func (pg *Storage) UpdateTransactionStatus(ctx context.Context, transactionID int, status string) error {
	const op = "postgres.UpdateTransactionStatus"

	query := `
        UPDATE transactions
        SET status = $1, processed_at = NOW()
        WHERE id = $2;
    `
	_, err := pg.db.Exec(ctx, query, status, transactionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Storage) GetTransactionsByClientID(ctx context.Context, clientID int) ([]models.Transaction, error) {
	const op = "postgres.GetTransactionsByClientID"

	query := `
        SELECT id, from_client_id, to_client_id, amount, status, created_at, processed_at
        FROM transactions
        WHERE from_client_id = $1 OR to_client_id = $1;
    `
	rows, err := pg.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.FromClientID,
			&t.ToClientID,
			&t.Amount,
			&t.Status,
			&t.CreatedAt,
			&t.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return transactions, nil
}

func (pg *Storage) GetTransaction(ctx context.Context, transactionID int) (models.Transaction, error) {
	const op = "postgres.GetTransactionsByClientID"

	var transaction models.Transaction

	query := `
        SELECT id, from_client_id, to_client_id, amount, status, created_at, processed_at
        FROM transactions
        WHERE id = $1
    `

	var processedTime *time.Time
	err := pg.db.QueryRow(ctx, query, transactionID).Scan(
		&transaction.ID,
		&transaction.FromClientID,
		&transaction.ToClientID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&processedTime,
	)
	if processedTime != nil {
		transaction.ProcessedAt = *processedTime
	}

	if err == pgx.ErrNoRows {
		return transaction, fmt.Errorf("client with id %d not found", transactionID)
	} else if err != nil {
		return transaction, fmt.Errorf("%s: %w", op, err)
	}

	return transaction, nil
}

func (pg *Storage) GetTransactionsByStatus(ctx context.Context, status string) ([]models.Transaction, error) {
	const op = "postgres.GetTransactionByStatus"

	query := `
        SELECT id, from_client_id, to_client_id, amount, status, created_at, processed_at
        FROM transactions
        WHERE status = $1;
    `
	rows, err := pg.db.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var processedTime *time.Time
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.FromClientID,
			&t.ToClientID,
			&t.Amount,
			&t.Status,
			&t.CreatedAt,
			&processedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if processedTime != nil {
			t.ProcessedAt = *processedTime
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return transactions, nil
}

func (pg *Storage) GetTransactionsByStatusAndClientID(ctx context.Context, clientID int, status string) ([]models.Transaction, error) {
	const op = "postgres.GetTransactionsByStatusAndClientID"

	query := `
        SELECT id, from_client_id, to_client_id, amount, status, created_at, processed_at
        FROM transactions
        WHERE (from_client_id = $1 OR to_client_id = $1) AND status = $2;
    `
	rows, err := pg.db.Query(ctx, query, clientID, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var processedTime *time.Time
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.FromClientID,
			&t.ToClientID,
			&t.Amount,
			&t.Status,
			&t.CreatedAt,
			&processedTime,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if processedTime != nil {
			t.ProcessedAt = *processedTime
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return transactions, nil
}
