// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"net/mail"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Define service errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type authService struct {
	repo   UserRepository
	hasher PasswordHasher
}

func (s *authService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewAuthService(repo UserRepository, hasher PasswordHasher) AuthService {
	return &authService{
		repo:   repo,
		hasher: hasher,
	}
}

// Email validation
func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmailFormat
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmailFormat
	}
	return nil
}

// Password validation
func validatePassword(pw string) error {
	if len(pw) < 8 {
		return ErrPasswordTooShort
	}

	var hasLetter, hasDigit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return ErrPasswordTooWeak
	}

	return nil
}

func (s *authService) SignUp(ctx context.Context, email, password string) error {
	// email/password validation
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
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
		if errors.Is(err, ErrEmailAlreadyExists) {
			return ErrEmailAlreadyExists
		}
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
