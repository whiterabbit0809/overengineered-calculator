// internal/auth/password.go
package auth

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) error
}

type bcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() PasswordHasher {
	return &bcryptPasswordHasher{}
}

func (b *bcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (b *bcryptPasswordHasher) CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
