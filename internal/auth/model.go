// internal/auth/model.go
package auth

import "time"

// Defines the User data model
type User struct {
	ID        string
	Email     string
	Password  string // hashed password
	CreatedAt time.Time
}
