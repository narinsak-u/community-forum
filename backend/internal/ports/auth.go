// Package ports defines the interfaces (contracts) for the application.
package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

// AuthService is an "Inbound Port" that defines the contract for authentication
// and authorization business logic.
type AuthService interface {
	Signup(ctx context.Context, username, email, password string) error
	Signin(ctx context.Context, login, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}
