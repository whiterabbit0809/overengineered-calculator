// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Define service errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (bool, error) // passed/failed
}

type authService struct {
	repo   UserRepository
	hasher PasswordHasher
}

func NewAuthService(repo UserRepository, hasher PasswordHasher) AuthService {
	return &authService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *authService) SignUp(ctx context.Context, email, password string) error {
	// TODO: add better email/password validation later
	if email == "" || password == "" {
		return errors.New("email and password required")
	}

	hash, err := s.hasher.HashPassword(password)
	if err != nil {
		return err
	}

	user := User{
		ID:        uuid.NewString(),
		Email:     email,
		Password:  hash,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *authService) Login(ctx context.Context, email, password string) (bool, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// login failed
			return false, nil
		}
		return false, err
	}

	if err := s.hasher.CheckPassword(user.Password, password); err != nil {
		// wrong password
		return false, nil
	}

	return true, nil
}
