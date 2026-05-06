package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

type CommentService interface {
	Create(ctx context.Context, slug string, content string, parentID *uint, authorID uint) (*domain.Comment, error)
	Delete(ctx context.Context, id uint, userID uint, userRole string) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id uint) error
	DeleteReplies(ctx context.Context, parentID uint) error
	GetByID(ctx context.Context, id uint) (*domain.Comment, error)
}
