package handlers_test

import (
	"context"
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var errMock = errors.New("mock error")

func init() {
	store := session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		CookieSecure:   false,
	})
	middleware.Store = store
}

type mockAuthService struct {
	signupFn  func(ctx context.Context, username, email, password string) error
	signinFn  func(ctx context.Context, login, password string) (*domain.User, error)
	getByIDFn func(ctx context.Context, id uint) (*domain.User, error)
}

func (m *mockAuthService) Signup(ctx context.Context, username, email, password string) error {
	return m.signupFn(ctx, username, email, password)
}

func (m *mockAuthService) Signin(ctx context.Context, login, password string) (*domain.User, error) {
	return m.signinFn(ctx, login, password)
}

func (m *mockAuthService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}

type mockUserService struct {
	getUserProfileFn func(ctx context.Context, username string) (*domain.User, error)
	updateProfileFn  func(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error)
}

func (m *mockUserService) GetUserProfile(ctx context.Context, username string) (*domain.User, error) {
	return m.getUserProfileFn(ctx, username)
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error) {
	return m.updateProfileFn(ctx, userID, updates)
}

var _ ports.AuthService = (*mockAuthService)(nil)
var _ ports.UserService = (*mockUserService)(nil)
