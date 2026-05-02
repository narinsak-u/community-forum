package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, username, email, password string) error
	Signin(ctx context.Context, login, password string) (*domain.User, error)
}
