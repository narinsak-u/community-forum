package db

import (
	"context"
	"strings"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMTagRepository struct {
	db *gorm.DB
}

func NewGORMTagRepository(db *gorm.DB) *GORMTagRepository {
	return &GORMTagRepository{db: db}
}

func (r *GORMTagRepository) ListTags(ctx context.Context) ([]domain.Tag, error) {
	var modelsTags []models.Tag
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&modelsTags).Error; err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(modelsTags))
	for i, t := range modelsTags {
		tags[i] = domain.Tag{
			ID:    t.ID,
			Name:  t.Name,
			Color: t.Color,
		}
	}
	return tags, nil
}

func (r *GORMTagRepository) CreateTag(ctx context.Context, tag *domain.Tag) error {
	m := &models.Tag{
		Name:  tag.Name,
		Color: tag.Color,
	}

	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	tag.ID = m.ID
	return nil
}

func (r *GORMTagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	var m models.Tag
	if err := r.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(name)).First(&m).Error; err != nil {
		return nil, err
	}

	return &domain.Tag{
		ID:    m.ID,
		Name:  m.Name,
		Color: m.Color,
	}, nil
}
