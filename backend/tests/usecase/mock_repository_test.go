package usecase_test

import (
	"context"
	"errors"
	"time"

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

type mockChatRepo struct {
	messages []domain.ChatMessage
	nextID   uint
	createFn func(msg *domain.ChatMessage) error
}

func newMockChatRepo() *mockChatRepo {
	return &mockChatRepo{
		messages: make([]domain.ChatMessage, 0),
		nextID:   1,
	}
}

func (m *mockChatRepo) Create(ctx context.Context, msg *domain.ChatMessage) error {
	if m.createFn != nil {
		return m.createFn(msg)
	}
	msg.ID = m.nextID
	m.nextID++
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()
	m.messages = append(m.messages, *msg)
	return nil
}

func (m *mockChatRepo) GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	n := len(m.messages)
	if n == 0 {
		return []domain.ChatMessage{}, nil
	}
	start := n - limit
	if start < 0 {
		start = 0
	}
	result := make([]domain.ChatMessage, 0, n-start)
	for i := start; i < n; i++ {
		result = append(result, m.messages[i])
	}
	return result, nil
}

func (m *mockChatRepo) GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	result := make([]domain.ChatMessage, 0)
	for _, msg := range m.messages {
		if msg.ID < beforeID {
			result = append(result, msg)
		}
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

var errInternal = errors.New("internal error")
