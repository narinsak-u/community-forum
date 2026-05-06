package db

import (
	"context"

	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMVoteRepository struct {
	db *gorm.DB
}

func NewGORMVoteRepository(db *gorm.DB) *GORMVoteRepository {
	return &GORMVoteRepository{db: db}
}

func (r *GORMVoteRepository) VoteThread(ctx context.Context, threadID uint, userID uint, value int8) error {
	if value == 0 {
		return r.db.WithContext(ctx).Where("user_id = ? AND thread_id = ?", userID, threadID).Delete(&models.Vote{}).Error
	}

	var existingVote models.Vote
	result := r.db.WithContext(ctx).Where("user_id = ? AND thread_id = ?", userID, threadID).First(&existingVote)

	if result.RowsAffected > 0 {
		return r.db.WithContext(ctx).Model(&existingVote).Update("value", value).Error
	}

	vote := models.Vote{
		UserID:   userID,
		ThreadID: &threadID,
		Value:    value,
	}
	return r.db.WithContext(ctx).Create(&vote).Error
}

func (r *GORMVoteRepository) VoteComment(ctx context.Context, commentID uint, userID uint, value int8) error {
	if value == 0 {
		return r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.Vote{}).Error
	}

	var existingVote models.Vote
	result := r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingVote)

	if result.RowsAffected > 0 {
		return r.db.WithContext(ctx).Model(&existingVote).Update("value", value).Error
	}

	vote := models.Vote{
		UserID:    userID,
		CommentID: &commentID,
		Value:     value,
	}
	return r.db.WithContext(ctx).Create(&vote).Error
}

func (r *GORMVoteRepository) GetThreadVotes(ctx context.Context, threadID uint) (int64, int64, error) {
	var upvotes, downvotes int64
	r.db.WithContext(ctx).Model(&models.Vote{}).Where("thread_id = ? AND value = 1", threadID).Count(&upvotes)
	r.db.WithContext(ctx).Model(&models.Vote{}).Where("thread_id = ? AND value = -1", threadID).Count(&downvotes)
	return upvotes, downvotes, nil
}

func (r *GORMVoteRepository) GetCommentVotes(ctx context.Context, commentID uint) (int64, int64, error) {
	var upvotes, downvotes int64
	r.db.WithContext(ctx).Model(&models.Vote{}).Where("comment_id = ? AND value = 1", commentID).Count(&upvotes)
	r.db.WithContext(ctx).Model(&models.Vote{}).Where("comment_id = ? AND value = -1", commentID).Count(&downvotes)
	return upvotes, downvotes, nil
}
