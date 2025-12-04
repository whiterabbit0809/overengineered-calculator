// internal/auth/model.go
package auth

import (
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
