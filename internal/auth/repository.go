// internal/auth/repository.go
package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Define repository errors
var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password, created_at)
         VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, user.Password, user.CreatedAt,
	)
	if err != nil {
		// Check for unique constraint violation on email
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password, created_at FROM users WHERE email = $1`,
		email,
	)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

// isUniqueViolation detects Postgres unique-constraint errors (email already exists).
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// 23505 = unique_violation
		return string(pqErr.Code) == "23505"
	}
	return false
}
