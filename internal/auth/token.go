package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims defines what we store in the JWT.
type TokenClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenService is an interface so we can swap implementation or mock in tests.
type TokenService interface {
	GenerateToken(user *User) (string, error)
	ParseToken(tokenStr string) (*TokenClaims, error)
}

// jwtTokenService is our concrete implementation using HS256.
type jwtTokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTTokenService(secret, issuer string, ttl time.Duration) TokenService {
	return &jwtTokenService{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (s *jwtTokenService) GenerateToken(user *User) (string, error) {
	now := time.Now().UTC()

	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtTokenService) ParseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what we expect
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %T", t.Method)
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
