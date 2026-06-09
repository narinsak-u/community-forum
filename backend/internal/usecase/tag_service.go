package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
)

type TagService struct {
	repo ports.TagRepository
}

func NewTagService(repo ports.TagRepository) *TagService {
	return &TagService{repo: repo}
}

var hexColorRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (s *TagService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.ListTags(ctx)
}

func (s *TagService) CreateTag(ctx context.Context, name, color, userRole string) (*domain.Tag, error) {
	if userRole != "admin" {
		return nil, errors.New("Admin access required")
	}

	name = strings.TrimSpace(name)
	color = strings.TrimSpace(color)

	if len(name) < 3 || len(name) > 50 {
		return nil, errors.New("Tag name must be between 3 and 50 characters")
	}

	if color != "" && !hexColorRegex.MatchString(color) {
		return nil, errors.New("Invalid color format, must be a valid hex color (e.g. #6366f1)")
	}

	if color == "" {
		color = "#6366f1"
	}

	if _, err := s.repo.GetByName(ctx, name); err == nil {
		return nil, errors.New("Tag with this name already exists")
	}

	tag := &domain.Tag{
		Name:  name,
		Color: color,
	}

	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, errors.New("Failed to create tag")
	}

	return tag, nil
}
