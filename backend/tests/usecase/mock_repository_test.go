package usecase_test

import (
	"context"
	"errors"

	"community-forum/backend/internal/domain"
)

type mockUserRepo struct {
	usersByID       map[uint]*domain.User
	usersByUsername map[string]*domain.User
	usersByEmail    map[string]*domain.User
	createFn        func(user *domain.User) error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:       make(map[uint]*domain.User),
		usersByUsername: make(map[string]*domain.User),
		usersByEmail:    make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(user)
	}
	m.usersByID[user.ID] = user
	m.usersByUsername[user.Username] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user, ok := m.usersByID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, ok := m.usersByUsername[username]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) List(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.usersByID))
	for _, u := range m.usersByID {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	existing, ok := m.usersByID[user.ID]
	if !ok {
		return domain.ErrNotFound
	}
	*existing = *user
	return nil
}

func (m *mockUserRepo) addUser(user *domain.User) {
	m.usersByID[user.ID] = user
	m.usersByUsername[user.Username] = user
	m.usersByEmail[user.Email] = user
}

var errInternal = errors.New("internal error")
