// Package ports defines the interfaces (contracts) for the application.
// In Hexagonal Architecture, ports allow the core logic to interact with external
// components (like databases or APIs) without knowing their implementation details.
package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

// UserRepository is an "Outbound Port". It defines the contract for how the
// application's core logic can interact with a data store to manage users.
type UserRepository interface {
	// context.Context is used to handle timeouts, cancellations, and request-scoped values.
	// It is a standard practice in Go to pass context as the first argument to methods
	// that perform I/O operations.
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// UserService is an "Inbound Port". It defines the contract for user-related
// business operations that can be triggered by external actors (like HTTP handlers).
type UserService interface {
	GetUserProfile(ctx context.Context, username string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error)
}
