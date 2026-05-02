package db

import (
	"context"
	"encoding/json"
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

func (r *GORMUserRepository) Create(ctx context.Context, u *domain.User) error {
	m := toModel(u)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GORMUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return fromModel(&m), nil
}

func (r *GORMUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return fromModel(&m), nil
}

func (r *GORMUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return fromModel(&m), nil
}

func (r *GORMUserRepository) Update(ctx context.Context, u *domain.User) error {
	m := toModel(u)
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", m.ID).Updates(m).Error
}

func toModel(u *domain.User) *models.User {
	stacks, _ := json.Marshal(u.Stacks)
	return &models.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Stacks:    string(stacks),
		Role:      u.Role,
	}
}

func fromModel(m *models.User) *domain.User {
	var stacks []string
	if m.Stacks != "" {
		_ = json.Unmarshal([]byte(m.Stacks), &stacks)
	}
	return &domain.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Avatar:    m.Avatar,
		Bio:       m.Bio,
		Stacks:    stacks,
		Role:      m.Role,
	}
}
