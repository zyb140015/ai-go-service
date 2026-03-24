package domain

import "time"

// AuthClaims represents the verified identity extracted from a bearer token.
type AuthClaims struct {
	UserID int64
	Email  string
	Expiry time.Time
}

// User represents an authenticated application user.
type User struct {
	ID           int64
	Email        string
	DisplayName  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
