package db

import (
	"context"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMCommentRepository struct {
	db *gorm.DB
}

func NewGORMCommentRepository(db *gorm.DB) *GORMCommentRepository {
	return &GORMCommentRepository{db: db}
}

func (r *GORMCommentRepository) Create(ctx context.Context, c *domain.Comment) error {
	m := &models.Comment{
		Content:  c.Content,
		ThreadID: c.ThreadID,
		AuthorID: c.AuthorID,
		ParentID: c.ParentID,
	}

	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Preload("Author").First(m, m.ID).Error; err != nil {
		return err
	}

	domainComment := commentFromModel(m, r.db)
	*c = *domainComment
	return nil
}

func (r *GORMCommentRepository) GetByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var m models.Comment
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return commentFromModel(&m, r.db), nil
}

func (r *GORMCommentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}

func (r *GORMCommentRepository) DeleteReplies(ctx context.Context, parentID uint) error {
	return r.db.WithContext(ctx).Where("parent_id = ?", parentID).Delete(&models.Comment{}).Error
}
