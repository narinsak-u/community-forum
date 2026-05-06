package usecase

import (
	"context"
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

type ThreadService struct {
	repo ports.ThreadRepository
}

func NewThreadService(repo ports.ThreadRepository) *ThreadService {
	return &ThreadService{repo: repo}
}

func (s *ThreadService) Create(ctx context.Context, title, content, status string, tags []string, authorID uint) (*domain.Thread, error) {
	if len(title) < 5 || len(title) > 255 {
		return nil, errors.New("Title must be between 5 and 255 characters")
	}

	if len(content) < 10 || len(content) > 50000 {
		return nil, errors.New("Content must be between 10 and 50000 characters")
	}

	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "published" {
		return nil, errors.New("Status must be draft or published")
	}

	slug, err := s.repo.GenerateUniqueSlug(ctx, title)
	if err != nil {
		return nil, errors.New("Failed to generate slug")
	}

	thread := &domain.Thread{
		Title:    title,
		Slug:     slug,
		Content:  content,
		Status:   status,
		AuthorID: authorID,
	}

	err = s.repo.Create(ctx, thread, tags)
	if err != nil {
		return nil, errors.New("Failed to create thread")
	}

	return thread, nil
}

func (s *ThreadService) List(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return s.repo.List(ctx, page, pageSize, sort)
}

func (s *ThreadService) ListByUser(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return s.repo.ListByUser(ctx, username, page, pageSize)
}

func (s *ThreadService) GetFeatured(ctx context.Context) (*domain.Thread, error) {
	return s.repo.GetFeatured(ctx)
}

func (s *ThreadService) GetTrending(ctx context.Context) ([]domain.Thread, error) {
	return s.repo.GetTrending(ctx)
}

func (s *ThreadService) GetBySlug(ctx context.Context, slug string) (*domain.Thread, error) {
	thread, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	_ = s.repo.IncrementViewCount(ctx, thread.ID)
	thread.ViewCount++

	return thread, nil
}

func (s *ThreadService) Update(ctx context.Context, slug string, authorID uint, title, content, status *string, tags []string) (*domain.Thread, error) {
	thread, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("Thread not found")
	}

	if thread.AuthorID != authorID {
		return nil, errors.New("You do not have permission to update this thread")
	}

	if title != nil {
		if len(*title) < 5 || len(*title) > 255 {
			return nil, errors.New("Title must be between 5 and 255 characters")
		}
		newSlug, err := s.repo.GenerateUniqueSlug(ctx, *title)
		if err != nil {
			return nil, errors.New("Failed to regenerate slug")
		}
		thread.Title = *title
		thread.Slug = newSlug
	}

	if content != nil {
		if len(*content) < 10 || len(*content) > 50000 {
			return nil, errors.New("Content must be between 10 and 50000 characters")
		}
		thread.Content = *content
	}

	if status != nil {
		if *status != "draft" && *status != "published" {
			return nil, errors.New("Status must be draft or published")
		}
		thread.Status = *status
	}

	err = s.repo.Update(ctx, thread, tags)
	if err != nil {
		return nil, errors.New("Failed to update thread")
	}

	return thread, nil
}

func (s *ThreadService) Delete(ctx context.Context, slug string, authorID uint) error {
	thread, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return errors.New("Thread not found")
	}

	if thread.AuthorID != authorID {
		return errors.New("You do not have permission to delete this thread")
	}

	return s.repo.Delete(ctx, thread)
}
