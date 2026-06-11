package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

// ChatRepository is an "Outbound Port" for chat message persistence.
type ChatRepository interface {
	Create(ctx context.Context, msg *domain.ChatMessage) error
	GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}

// ChatService is an "Inbound Port" for chat business operations.
type ChatService interface {
	SendMessage(ctx context.Context, authorID uint, content string) (*domain.ChatMessage, error)
	GetRecentMessages(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetMessagesBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}
