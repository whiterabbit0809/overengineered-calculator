// internal/auth/model.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Defines the User data model
type User struct {
	ID        string
	Email     string
	Password  string // hashed password
	CreatedAt time.Time
}

// TokenClaims defines what we store in the JWT.
type TokenClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Custom error definitions
var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak    = errors.New("password must contain at least one letter and one digit")
)

// Response models for handlers
type signUpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Handler struct {
	service      AuthService
	tokenService TokenService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"`          // "passed" or "failed"
	Token  string `json:"token,omitempty"` // JWT token if passed
}
