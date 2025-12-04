// internal/storage/postgres.go
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
`
const createHistoryTable = `
CREATE TABLE IF NOT EXISTS calc_history (
    id          SERIAL PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users(id),
    expression  TEXT        NOT NULL,
    result      DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

`

func NewPostgresDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback to individual env vars (for docker-compose)
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		pass := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "overengineered_calc")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, pass, host, port, name, sslmode)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Retry ping a few times – Postgres might not be ready yet in Docker.
	if err := pingWithRetry(db, 10, time.Second); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	// Auto-create users table
	if _, err := db.Exec(createUsersTableQuery); err != nil {
		return nil, fmt.Errorf("create users table: %w", err)
	}
	// Auto-create calc_history table
	if _, err := db.Exec(createHistoryTable); err != nil {
		return nil, fmt.Errorf("create calc_history table: %w", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func pingWithRetry(db *sql.DB, attempts int, delay time.Duration) error {
	var err error
	for attempts > 0 {
		err = db.Ping()
		if err == nil {
			return nil
		}
		attempts--
		if attempts == 0 {
			break
		}
		time.Sleep(delay)
	}
	return err
}
