package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPostgres(dbURL string) error {
	var err error

	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("create postgres pool: %w", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		DB.Close()
		DB = nil
		return fmt.Errorf("connect to postgres: %w", err)
	}

	return nil
}

func CreateTables() error {
	_, err := DB.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		`,
	)

	if err != nil {
		return fmt.Errorf("database: create 'users' table error: %w", err)
	}

	_, err = DB.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS chats (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
	)

	if err != nil {
		return fmt.Errorf("database: create 'chats' table error: %w", err)
	}

	_, err = DB.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS messages (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
	)

	if err != nil {
		return fmt.Errorf("database: create 'messages' table error: %w", err)
	}

	_, err = DB.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token TEXT NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL
		);
		`,
	)

	if err != nil {
		return fmt.Errorf("database: create 'refresh_tokens' table error: %w", err)
	}

	_, err = DB.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS api_keys (
		    id BIGSERIAL PRIMARY KEY,
		    title TEXT NOT NULL,
		    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		    key_hash TEXT NOT NULL UNIQUE,
		    expires_at TIMESTAMPTZ DEFAULT NULL,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		    UNIQUE (title, user_id)
		);
		`,
	)

	if err != nil {
		return fmt.Errorf("database: create 'api_keys' table error: %w", err)
	}

	return nil
}
