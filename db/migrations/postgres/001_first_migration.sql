-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance NUMERIC(15, 2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    to_client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS transaction_queue (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    queued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'processed', 'failed'))
);

CREATE TABLE IF NOT EXISTS transaction_history (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    from_client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    to_client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    balance_before NUMERIC(15, 2) NOT NULL,
    balance_after NUMERIC(15, 2) NOT NULL,
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS transaction_history;
DROP TABLE IF EXISTS transaction_queue;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS clients;

-- +goose StatementEnd
