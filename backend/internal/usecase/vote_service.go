package usecase

import (
	"context"
	"errors"

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
		return 0, 0, errors.New("Value must be -1, 0, or 1")
	}

	thread, err := s.threadRepo.GetBySlug(ctx, slug)
	if err != nil {
		return 0, 0, errors.New("Thread not found")
	}

	if err := s.repo.VoteThread(ctx, thread.ID, userID, value); err != nil {
		return 0, 0, err
	}

	return s.repo.GetThreadVotes(ctx, thread.ID)
}

func (s *VoteService) VoteComment(ctx context.Context, commentID uint, userID uint, value int8) (int64, int64, error) {
	if value != -1 && value != 0 && value != 1 {
		return 0, 0, errors.New("Value must be -1, 0, or 1")
	}

	_, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return 0, 0, errors.New("Comment not found")
	}

	if err := s.repo.VoteComment(ctx, commentID, userID, value); err != nil {
		return 0, 0, err
	}

	return s.repo.GetCommentVotes(ctx, commentID)
}
