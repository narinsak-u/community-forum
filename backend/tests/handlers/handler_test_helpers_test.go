package handlers_test

import (
	"context"
	"errors"
	"time"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"
)

var errMock = errors.New("mock error")

func setupTestSessionManager() *middleware.SessionManager {
	return middleware.NewSessionManager("test-secret", 72*time.Hour)
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
	getUserByIDFn    func(ctx context.Context, id uint) (*domain.User, error)
	listUsersFn      func(ctx context.Context) ([]*domain.User, error)
	updateProfileFn  func(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error)
}

func (m *mockUserService) GetUserProfile(ctx context.Context, username string) (*domain.User, error) {
	return m.getUserProfileFn(ctx, username)
}

func (m *mockUserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return m.getUserByIDFn(ctx, id)
}

func (m *mockUserService) ListUsers(ctx context.Context) ([]*domain.User, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(ctx)
	}
	return nil, nil
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error) {
	return m.updateProfileFn(ctx, userID, updates)
}

var _ ports.AuthService = (*mockAuthService)(nil)
var _ ports.UserService = (*mockUserService)(nil)

type mockThreadService struct {
	createFn      func(ctx context.Context, title, content, status string, tags []string, authorID uint) (*domain.Thread, error)
	listFn        func(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error)
	listByUserFn  func(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error)
	getFeaturedFn func(ctx context.Context) (*domain.Thread, error)
	getTrendingFn func(ctx context.Context) ([]domain.Thread, error)
	getBySlugFn   func(ctx context.Context, slug string) (*domain.Thread, error)
	updateFn      func(ctx context.Context, slug string, authorID uint, title, content, status *string, tags []string) (*domain.Thread, error)
	deleteFn      func(ctx context.Context, slug string, authorID uint) error
}

func (m *mockThreadService) Create(ctx context.Context, title, content, status string, tags []string, authorID uint) (*domain.Thread, error) {
	return m.createFn(ctx, title, content, status, tags, authorID)
}

func (m *mockThreadService) List(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error) {
	return m.listFn(ctx, page, pageSize, sort)
}

func (m *mockThreadService) ListByUser(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error) {
	return m.listByUserFn(ctx, username, page, pageSize)
}

func (m *mockThreadService) GetFeatured(ctx context.Context) (*domain.Thread, error) {
	return m.getFeaturedFn(ctx)
}

func (m *mockThreadService) GetTrending(ctx context.Context) ([]domain.Thread, error) {
	return m.getTrendingFn(ctx)
}

func (m *mockThreadService) GetBySlug(ctx context.Context, slug string) (*domain.Thread, error) {
	return m.getBySlugFn(ctx, slug)
}

func (m *mockThreadService) Update(ctx context.Context, slug string, authorID uint, title, content, status *string, tags []string) (*domain.Thread, error) {
	return m.updateFn(ctx, slug, authorID, title, content, status, tags)
}

func (m *mockThreadService) Delete(ctx context.Context, slug string, authorID uint) error {
	return m.deleteFn(ctx, slug, authorID)
}

var _ ports.ThreadService = (*mockThreadService)(nil)
