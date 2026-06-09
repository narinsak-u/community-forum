package usecase_test

import (
	"context"
	"testing"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetUserProfile_Success(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Bio:      "Hello world",
		Stacks:   []string{"Go", "React"},
	})
	svc := usecase.NewUserService(mock)

	user, err := svc.GetUserProfile(context.Background(), "johndoe")
	require.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "Hello world", user.Bio)
	assert.Equal(t, []string{"Go", "React"}, user.Stacks)
}

func TestUserService_GetUserProfile_NotFound(t *testing.T) {
	svc := usecase.NewUserService(newMockUserRepo())

	_, err := svc.GetUserProfile(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Bio:      "Original bio",
		Avatar:   "original.jpg",
		Stacks:   []string{"Go"},
	})
	svc := usecase.NewUserService(mock)

	updates := &domain.User{
		Bio:    "Updated bio",
		Avatar: "new-avatar.jpg",
		Stacks: []string{"Go", "React", "TypeScript"},
	}

	user, err := svc.UpdateProfile(context.Background(), 1, updates)
	require.NoError(t, err)
	assert.Equal(t, "Updated bio", user.Bio)
	assert.Equal(t, "new-avatar.jpg", user.Avatar)
	assert.Equal(t, []string{"Go", "React", "TypeScript"}, user.Stacks)
}

func TestUserService_UpdateProfile_PartialUpdate(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Bio:      "Original bio",
		Avatar:   "original.jpg",
		Stacks:   []string{"Go"},
	})
	svc := usecase.NewUserService(mock)

	updates := &domain.User{Bio: "Just updating my bio"}

	user, err := svc.UpdateProfile(context.Background(), 1, updates)
	require.NoError(t, err)
	assert.Equal(t, "Just updating my bio", user.Bio)
	assert.Equal(t, "original.jpg", user.Avatar)
	assert.Equal(t, []string{"Go"}, user.Stacks)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	svc := usecase.NewUserService(newMockUserRepo())

	updates := &domain.User{Bio: "New bio"}
	_, err := svc.UpdateProfile(context.Background(), 999, updates)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUserService_UpdateProfile_EmptyAvatarDoesNotOverwrite(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Avatar:   "original.jpg",
		Stacks:   []string{},
	})
	svc := usecase.NewUserService(mock)

	updates := &domain.User{}

	user, err := svc.UpdateProfile(context.Background(), 1, updates)
	require.NoError(t, err)
	assert.Equal(t, "original.jpg", user.Avatar)
}

func TestUserService_UpdateProfile_StacksNilDoesNotOverwrite(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Stacks:   []string{"Go", "React"},
	})
	svc := usecase.NewUserService(mock)

	user, err := svc.UpdateProfile(context.Background(), 1, &domain.User{})
	require.NoError(t, err)
	assert.Equal(t, []string{"Go", "React"}, user.Stacks)
}

func TestUserService_UpdateProfile_EmptySticksOverwrites(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Stacks:   []string{"Go", "React"},
	})
	svc := usecase.NewUserService(mock)

	user, err := svc.UpdateProfile(context.Background(), 1, &domain.User{Stacks: []string{}})
	require.NoError(t, err)
	assert.Equal(t, []string{}, user.Stacks)
}
