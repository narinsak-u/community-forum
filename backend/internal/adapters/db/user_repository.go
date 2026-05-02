// Package db implements the database adapters for the application.
// In Hexagonal Architecture, this is an "Outbound Adapter". It implements 
// an Outbound Port (UserRepository) using a specific technology (GORM/Postgres).
package db

import (
	"context"
	"encoding/json"
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

// GORMUserRepository implements ports.UserRepository.
type GORMUserRepository struct {
	// db is the GORM database connection.
	db *gorm.DB
}

// NewGORMUserRepository creates a new instance of GORMUserRepository.
func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

// Create persists a new user to the database.
func (r *GORMUserRepository) Create(ctx context.Context, u *domain.User) error {
	// Step 1: Map the domain User entity to the database models.User.
	m, err := toModel(u)
	if err != nil {
		return err
	}
	// Step 2: Use GORM to insert the model into the database.
	// WithContext ensures that the database operation respects the request context.
	return r.db.WithContext(ctx).Create(m).Error
}

// GetByID retrieves a user from the database by their ID.
func (r *GORMUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var m models.User
	// Step 1: Use GORM's First method to find the user by primary key.
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		// Step 2: If the record is not found, map GORM's error to a domain-specific error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	// Step 3: Map the database model back to the domain entity.
	return fromModel(&m)
}

// GetByUsername retrieves a user from the database by their username.
func (r *GORMUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return fromModel(&m)
}

// GetByEmail retrieves a user from the database by their email address.
func (r *GORMUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return fromModel(&m)
}

// Update modifies an existing user in the database.
func (r *GORMUserRepository) Update(ctx context.Context, u *domain.User) error {
	// Step 1: Map domain entity to database model.
	m, err := toModel(u)
	if err != nil {
		return err
	}
	// Step 2: Perform the update operation.
	// We explicitly select the fields to update and omit "CreatedAt" to prevent accidental changes.
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", m.ID).
		Select("*").
		Omit("CreatedAt").
		Updates(m)

	if result.Error != nil {
		return result.Error
	}
	// Step 3: Check if any row was actually updated.
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// toModel is a private helper that maps a domain.User to a models.User.
// This is necessary because the domain layer should not know about database tags or formats.
func toModel(u *domain.User) (*models.User, error) {
	// The Stacks slice is stored as a JSON string in the database.
	stacks, err := json.Marshal(u.Stacks)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

// fromModel is a private helper that maps a models.User back to a domain.User.
func fromModel(m *models.User) (*domain.User, error) {
	var stacks []string
	if m.Stacks != "" {
		// Parse the JSON string from the database back into a Go slice.
		if err := json.Unmarshal([]byte(m.Stacks), &stacks); err != nil {
			return nil, err
		}
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
	}, nil
}
