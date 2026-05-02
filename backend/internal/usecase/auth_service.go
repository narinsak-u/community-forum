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

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailRegistered    = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

type AuthService struct {
	repo ports.UserRepository
}

func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, username, email, password string) error {
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

	_, err := s.repo.GetByUsername(ctx, username)
	if err == nil {
		return ErrUsernameTaken
	} else if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	_, err = s.repo.GetByEmail(ctx, email)
	if err == nil {
		return ErrEmailRegistered
	} else if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	return s.repo.Create(ctx, user)
}

func (s *AuthService) Signin(ctx context.Context, login, password string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			user, err = s.repo.GetByEmail(ctx, login)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					return nil, ErrInvalidCredentials
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
