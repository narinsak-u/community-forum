package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type UserService interface {
	GetUserProfile(ctx context.Context, username string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error)
}
