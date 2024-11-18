package postgres

import (
	"context"
	"fmt"

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

func (pg *Storage) AddTransaction(ctx context.Context, transaction models.Transaction) (int, error) {
	const op = "postgres.AddTransaction"

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

	err := pg.db.QueryRow(ctx, query, transactionID).Scan(
		&transaction.ID,
		&transaction.FromClientID,
		&transaction.ToClientID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.ProcessedAt,
	)

	if err == pgx.ErrNoRows {
		return transaction, fmt.Errorf("client with id %d not found", transactionID)
	} else if err != nil {
		return transaction, fmt.Errorf("%s: %w", op, err)
	}

	return transaction, nil
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
