package usecase

import (
	"context"
	"errors"
	"strings"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

const DefaultChatLimit = 15

type ChatService struct {
	repo ports.ChatRepository
}

func NewChatService(repo ports.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SendMessage(ctx context.Context, authorID uint, content string) (*domain.ChatMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, errors.New("message content too long")
	}

	msg := &domain.ChatMessage{
		Content:  content,
		AuthorID: authorID,
	}
	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ChatService) GetRecentMessages(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = DefaultChatLimit
	}
	return s.repo.GetRecent(ctx, limit)
}

func (s *ChatService) GetMessagesBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = DefaultChatLimit
	}
	return s.repo.GetBefore(ctx, beforeID, limit)
}
