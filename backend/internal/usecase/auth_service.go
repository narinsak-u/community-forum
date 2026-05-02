// Package usecase implements the core business logic of the application.
// In Hexagonal Architecture, this layer (also called the Service layer) 
// implements the Inbound Ports and coordinates interactions with Outbound Ports.
package usecase

import (
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
	"context"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Pre-defining errors for common business rule violations.
var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Regular expressions for input validation.
var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// AuthService implements the ports.AuthService interface.
type AuthService struct {
	// Dependency Injection: AuthService depends on UserRepository to persist data.
	// It uses the interface type (ports.UserRepository) to remain decoupled 
	// from the actual database implementation.
	repo ports.UserRepository
}

// NewAuthService is a constructor function that returns a new instance of AuthService.
func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Signup handles the registration of a new user.
func (s *AuthService) Signup(ctx context.Context, username, email, password string) error {
	// Step 1: Sanitize and validate inputs.
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if len(username) < 3 || len(username) > 30 || !usernameRegex.MatchString(username) {
		return errors.New("username must be 3-30 characters and contain only letters, numbers, and underscores")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Step 2: Check if the username is already taken.
	_, err := s.repo.GetByUsername(ctx, username)
	if err == nil {
		return ErrUsernameTaken
	} else if !errors.Is(err, domain.ErrNotFound) {
		// If the error is something other than "not found", return it.
		return err
	}

	// Step 3: Check if the email is already registered.
	_, err = s.repo.GetByEmail(ctx, email)
	if err == nil {
		return ErrEmailRegistered
	} else if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	// Step 4: Hash the password for security.
	// We use bcrypt with a cost of 12, which is a balanced value for security and performance.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Step 5: Create a new domain User entity.
	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	// Step 6: Persist the new user via the repository.
	return s.repo.Create(ctx, user)
}

// Signin handles user authentication.
func (s *AuthService) Signin(ctx context.Context, login, password string) (*domain.User, error) {
	// Step 1: Attempt to find the user by username.
	user, err := s.repo.GetByUsername(ctx, login)
	if err != nil {
		// Step 2: If not found by username, attempt to find by email.
		if errors.Is(err, domain.ErrNotFound) {
			user, err = s.repo.GetByEmail(ctx, login)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					// We return a generic "invalid credentials" error for security.
					return nil, ErrInvalidCredentials
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Step 3: Compare the provided password with the stored hash.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Step 4: Return the authenticated user.
	return user, nil
}

// GetByID retrieves a user by their ID.
func (s *AuthService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
