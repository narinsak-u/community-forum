package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

// ThreadService defines the business logic methods for threads.
type ThreadService interface {
	Create(ctx context.Context, title, content, status string, tags []string, authorID uint) (*domain.Thread, error)
	List(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error)
	ListByUser(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error)
	GetFeatured(ctx context.Context) (*domain.Thread, error)
	GetTrending(ctx context.Context) ([]domain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Thread, error)
	Update(ctx context.Context, slug string, authorID uint, title, content, status *string, tags []string) (*domain.Thread, error)
	Delete(ctx context.Context, slug string, authorID uint) error
}

// ThreadRepository defines the data access methods for threads.
type ThreadRepository interface {
	Create(ctx context.Context, thread *domain.Thread, tagNames []string) error
	List(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error)
	ListByUser(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error)
	GetFeatured(ctx context.Context) (*domain.Thread, error)
	GetTrending(ctx context.Context) ([]domain.Thread, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Thread, error)
	Update(ctx context.Context, thread *domain.Thread, tagNames []string) error
	Delete(ctx context.Context, thread *domain.Thread) error
	IncrementViewCount(ctx context.Context, threadID uint) error
	GenerateUniqueSlug(ctx context.Context, title string) (string, error)
}
