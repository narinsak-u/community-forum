// Package domain contains the core business logic and entities of the application.
// In Hexagonal Architecture, this is the innermost layer, which is independent
// of any external frameworks, databases, or APIs.
package domain

import (
	"errors"
	"time"
)

// ErrNotFound is a domain-specific error used when a user cannot be found.
// Defining custom errors at the domain level allows other layers to check for
// specific error conditions without being coupled to database-specific errors.
var ErrNotFound = errors.New("user not found")

// User roles are defined as constants to avoid "magic strings" in the codebase.
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// User represents the core domain entity for a user in the system.
// This struct defines the data and behavior (methods) that are essential
// to the business, regardless of how it is stored in the database.
type User struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Stacks    []string  `json:"stacks,omitempty"`
	Role      string    `json:"role"`
}
