package usecase

import (
	"context"
	"fmt"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

var (
	ErrCommentNotFound  = fmt.Errorf("comment: %w", domain.ErrNotFound)
	ErrCommentForbidden = fmt.Errorf("comment: %w", domain.ErrForbidden)
)

type CommentService struct {
	repo       ports.CommentRepository
	threadRepo ports.ThreadRepository
}

func NewCommentService(repo ports.CommentRepository, threadRepo ports.ThreadRepository) *CommentService {
	return &CommentService{
		repo:       repo,
		threadRepo: threadRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, slug string, content string, parentID *uint, authorID uint) (*domain.Comment, error) {
	if len(content) < 2 || len(content) > 10000 {
		return nil, fmt.Errorf("%w: content must be between 2 and 10000 characters", domain.ErrInvalidInput)
	}

	thread, err := s.threadRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrThreadNotFound, slug)
	}

	if parentID != nil {
		parent, err := s.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("%w: %d", ErrCommentNotFound, *parentID)
		}

		if parent.ThreadID != thread.ID {
			return nil, fmt.Errorf("%w: parent comment does not belong to this thread", domain.ErrInvalidInput)
		}

		if parent.ParentID != nil {
			return nil, fmt.Errorf("%w: replies can only be 1 level deep", domain.ErrInvalidInput)
		}
	}

	comment := &domain.Comment{
		Content:  content,
		ThreadID: thread.ID,
		AuthorID: authorID,
		ParentID: parentID,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	return comment, nil
}

func (s *CommentService) Delete(ctx context.Context, id uint, userID uint, userRole string) error {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %d", ErrCommentNotFound, id)
	}

	if comment.AuthorID != userID && userRole != "admin" {
		return ErrCommentForbidden
	}

	_ = s.repo.DeleteReplies(ctx, comment.ID)

	if err := s.repo.Delete(ctx, comment.ID); err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	return nil
}
