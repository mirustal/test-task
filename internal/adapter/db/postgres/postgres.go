package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"

	"bank-service/pkg/config"
)

const migrationDir = "db/migrations/postgres"

type Storage struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Storage, error) {
	const op = "postgres.New"

	log = log.With(slog.String("op", op))

	url := dbStringConverter(cfg)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("failed to parse db config: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	poolConfig.MaxConns = int32(cfg.Postgres.MaxConn)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to ping db: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := applyMigrations(ctx, pool, migrationDir); err != nil {
		log.Error("can't migrate up: ", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:  pool,
		log: log,
	}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.db.Close()
	return nil
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	db := stdlib.OpenDBFromPool(pool)
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
