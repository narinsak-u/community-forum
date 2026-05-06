package ports

import "context"

type VoteService interface {
	VoteThread(ctx context.Context, slug string, userID uint, value int8) (int64, int64, error)
	VoteComment(ctx context.Context, commentID uint, userID uint, value int8) (int64, int64, error)
}

type VoteRepository interface {
	VoteThread(ctx context.Context, threadID uint, userID uint, value int8) error
	VoteComment(ctx context.Context, commentID uint, userID uint, value int8) error
	GetThreadVotes(ctx context.Context, threadID uint) (int64, int64, error)
	GetCommentVotes(ctx context.Context, commentID uint) (int64, int64, error)
}
