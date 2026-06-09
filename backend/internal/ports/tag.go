package ports

import (
	"context"

	"community-forum/backend/internal/domain"
)

type TagService interface {
	ListTags(ctx context.Context) ([]domain.Tag, error)
	CreateTag(ctx context.Context, name, color, userRole string) (*domain.Tag, error)
}

type TagRepository interface {
	ListTags(ctx context.Context) ([]domain.Tag, error)
	CreateTag(ctx context.Context, tag *domain.Tag) error
	GetByName(ctx context.Context, name string) (*domain.Tag, error)
}
