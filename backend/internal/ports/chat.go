package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

type ChatRepository interface {
	Save(ctx context.Context, msg *domain.ChatMessage) error
	GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}

type ChatService interface {
	SendMessage(ctx context.Context, authorID uint, content string) (*domain.ChatMessage, error)
	GetRecentMessages(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetMessagesBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}
