// internal/history/repository.go
package history

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	Create(ctx context.Context, entry *HistoryEntry) error
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]HistoryEntry, error)
	GetLatestResult(ctx context.Context, userID string) (float64, error)
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Create(ctx context.Context, e *HistoryEntry) error {
	row := r.DB.QueryRowContext(ctx,
		`INSERT INTO calc_history (user_id, expression, result)
         VALUES ($1, $2, $3)
         RETURNING id, created_at`,
		e.UserID, e.Expression, e.Result,
	)
	return row.Scan(&e.ID, &e.CreatedAt)
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]HistoryEntry, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, user_id, expression, result, created_at
         FROM calc_history
         WHERE user_id = $1
         ORDER BY created_at DESC
         LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []HistoryEntry
	for rows.Next() {
		var e HistoryEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Expression, &e.Result, &e.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, rows.Err()
}

// GetLatestResult returns the last stored result for a user, or 0 if none exist.
func (r *PostgresRepository) GetLatestResult(ctx context.Context, userID string) (float64, error) {
	row := r.DB.QueryRowContext(ctx, `
        SELECT result
        FROM calc_history
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `, userID)

	var result float64
	err := row.Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		// No history yet → start from 0
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return result, nil
}
