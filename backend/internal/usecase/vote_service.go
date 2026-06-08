package usecase

import (
	"context"
	"fmt"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

type VoteService struct {
	repo        ports.VoteRepository
	threadRepo  ports.ThreadRepository
	commentRepo ports.CommentRepository
}

func NewVoteService(repo ports.VoteRepository, threadRepo ports.ThreadRepository, commentRepo ports.CommentRepository) *VoteService {
	return &VoteService{
		repo:        repo,
		threadRepo:  threadRepo,
		commentRepo: commentRepo,
	}
}

func (s *VoteService) VoteThread(ctx context.Context, slug string, userID uint, value int8) (int64, int64, error) {
	if value != -1 && value != 0 && value != 1 {
		return 0, 0, fmt.Errorf("%w: value must be -1, 0, or 1", domain.ErrInvalidInput)
	}

	thread, err := s.threadRepo.GetBySlug(ctx, slug)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %s", ErrThreadNotFound, slug)
	}

	if err := s.repo.VoteThread(ctx, thread.ID, userID, value); err != nil {
		return 0, 0, fmt.Errorf("vote thread: %w", err)
	}

	return s.repo.GetThreadVotes(ctx, thread.ID)
}

func (s *VoteService) VoteComment(ctx context.Context, commentID uint, userID uint, value int8) (int64, int64, error) {
	if value != -1 && value != 0 && value != 1 {
		return 0, 0, fmt.Errorf("%w: value must be -1, 0, or 1", domain.ErrInvalidInput)
	}

	_, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %d", ErrCommentNotFound, commentID)
	}

	if err := s.repo.VoteComment(ctx, commentID, userID, value); err != nil {
		return 0, 0, fmt.Errorf("vote comment: %w", err)
	}

	return s.repo.GetCommentVotes(ctx, commentID)
}
