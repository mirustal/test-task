package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	"bank-service/internal/config"
)

const migrationDir = "db/migrations/postgres"

type Storage struct {
	db  *pgx.Conn
	log *slog.Logger
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Storage, error) {
	const op = "postgres.New"

	log = log.With(slog.String("op", op))

	url := dbStringConverter(cfg)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Error("can`t connect db: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Error("failed ping db: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := applyMigrations(ctx, conn, migrationDir); err != nil {
		log.Error("can't migrate up: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:  conn,
		log: log,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close(ctx)
}

func applyMigrations(ctx context.Context, conn *pgx.Conn, migrationsDir string) error {
	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	return goose.Up(db, migrationsDir)
}

func dbStringConverter(cfg *config.Config) string {
	urlConn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
	)

	return urlConn
}
