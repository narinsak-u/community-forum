// Package usecase implements the core business logic.
package usecase

import (
	"context"
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

// UserService implements the ports.UserService interface.
// It handles high-level business operations related to user profiles.
type UserService struct {
	// repo is an outbound port used to interact with the database.
	repo ports.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserProfile retrieves a user's domain entity by their username.
func (s *UserService) GetUserProfile(ctx context.Context, username string) (*domain.User, error) {
	// Directly call the repository to fetch the user.
	return s.repo.GetByUsername(ctx, username)
}

// UpdateProfile updates an existing user's profile information.
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, updates *domain.User) (*domain.User, error) {
	// Step 1: Fetch the current user from the repository.
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Step 2: Apply updates if the provided fields are not empty.
	// This "selective update" pattern ensures we don't overwrite data with empty values.
	if updates.Bio != "" {
		user.Bio = updates.Bio
	}
	if updates.Avatar != "" {
		user.Avatar = updates.Avatar
	}
	if updates.Stacks != nil {
		user.Stacks = updates.Stacks
	}

	// Step 3: Save the updated user entity back to the repository.
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}
