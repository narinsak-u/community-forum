package db

import (
	"context"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMChatRepository struct {
	db *gorm.DB
}

func NewGORMChatRepository(db *gorm.DB) *GORMChatRepository {
	return &GORMChatRepository{db: db}
}

func (r *GORMChatRepository) Create(ctx context.Context, msg *domain.ChatMessage) error {
	m := &models.ChatMessage{
		Content:  msg.Content,
		AuthorID: msg.AuthorID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Preload("Author").First(m, m.ID).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	msg.CreatedAt = m.CreatedAt
	msg.UpdatedAt = m.UpdatedAt
	msg.Author = domain.User{
		ID:       m.Author.ID,
		Username: m.Author.Username,
		Avatar:   m.Author.Avatar,
	}
	return nil
}

func (r *GORMChatRepository) GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	var ms []models.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Order("id DESC").
		Limit(limit).
		Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]domain.ChatMessage, 0, len(ms))
	for i := len(ms) - 1; i >= 0; i-- {
		result = append(result, *chatMessageFromModel(&ms[i]))
	}
	return result, nil
}

func (r *GORMChatRepository) GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	var ms []models.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Where("id < ?", beforeID).
		Order("id DESC").
		Limit(limit).
		Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]domain.ChatMessage, 0, len(ms))
	for i := len(ms) - 1; i >= 0; i-- {
		result = append(result, *chatMessageFromModel(&ms[i]))
	}
	return result, nil
}

func chatMessageFromModel(m *models.ChatMessage) *domain.ChatMessage {
	return &domain.ChatMessage{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Content:   m.Content,
		AuthorID:  m.AuthorID,
		Author: domain.User{
			ID:       m.Author.ID,
			Username: m.Author.Username,
			Avatar:   m.Author.Avatar,
		},
	}
}
