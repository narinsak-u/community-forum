package usecase

import (
	"context"
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
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
		return nil, errors.New("Content must be between 2 and 10000 characters")
	}

	thread, err := s.threadRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("Thread not found")
	}

	if parentID != nil {
		parent, err := s.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, errors.New("Parent comment not found")
		}

		if parent.ThreadID != thread.ID {
			return nil, errors.New("Parent comment does not belong to this thread")
		}

		if parent.ParentID != nil {
			return nil, errors.New("Replies can only be 1 level deep")
		}
	}

	comment := &domain.Comment{
		Content:  content,
		ThreadID: thread.ID,
		AuthorID: authorID,
		ParentID: parentID,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, errors.New("Failed to create comment")
	}

	return comment, nil
}

func (s *CommentService) Delete(ctx context.Context, id uint, userID uint, userRole string) error {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("Comment not found")
	}

	if comment.AuthorID != userID && userRole != "admin" {
		return errors.New("You do not have permission to delete this comment")
	}

	_ = s.repo.DeleteReplies(ctx, comment.ID)

	if err := s.repo.Delete(ctx, comment.ID); err != nil {
		return errors.New("Failed to delete comment")
	}

	return nil
}
