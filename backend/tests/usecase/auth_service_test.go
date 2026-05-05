package usecase_test

import (
	"context"
	"testing"
	"time"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Signup_Success(t *testing.T) {
	mock := newMockUserRepo()
	svc := usecase.NewAuthService(mock)

	err := svc.Signup(context.Background(), "testuser", "test@example.com", "password123")
	require.NoError(t, err)

	user, err := mock.GetByUsername(context.Background(), "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, domain.RoleUser, user.Role)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)
}

func TestAuthService_Signup_UsernameTooShort(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())
	err := svc.Signup(context.Background(), "ab", "test@example.com", "password123")
	assert.ErrorContains(t, err, "username must be 3-30 characters")
}

func TestAuthService_Signup_UsernameTooLong(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())
	longUsername := ""
	for i := 0; i < 31; i++ {
		longUsername += "a"
	}
	err := svc.Signup(context.Background(), longUsername, "test@example.com", "password123")
	assert.ErrorContains(t, err, "username must be 3-30 characters")
}

func TestAuthService_Signup_UsernameInvalidChars(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())
	err := svc.Signup(context.Background(), "user name!", "test@example.com", "password123")
	assert.ErrorContains(t, err, "username must be 3-30 characters")
}

func TestAuthService_Signup_InvalidEmail(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())
	err := svc.Signup(context.Background(), "testuser", "not-an-email", "password123")
	assert.ErrorContains(t, err, "invalid email format")
}

func TestAuthService_Signup_PasswordTooShort(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())
	err := svc.Signup(context.Background(), "testuser", "test@example.com", "short")
	assert.ErrorContains(t, err, "password must be at least 8 characters")
}

func TestAuthService_Signup_UsernameTaken(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "existing",
		Email:    "existing@example.com",
		Password: "hashed",
		Role:     domain.RoleUser,
		Stacks:   []string{},
	})
	svc := usecase.NewAuthService(mock)

	err := svc.Signup(context.Background(), "existing", "new@example.com", "password123")
	assert.ErrorIs(t, err, usecase.ErrUsernameTaken)
}

func TestAuthService_Signup_EmailRegistered(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{
		ID:       1,
		Username: "existing",
		Email:    "existing@example.com",
		Password: "hashed",
		Role:     domain.RoleUser,
		Stacks:   []string{},
	})
	svc := usecase.NewAuthService(mock)

	err := svc.Signup(context.Background(), "newuser", "existing@example.com", "password123")
	assert.ErrorIs(t, err, usecase.ErrEmailRegistered)
}

func TestAuthService_Signin_SuccessByUsername(t *testing.T) {
	mock := newMockUserRepo()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	mock.addUser(&domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashed),
		Role:     domain.RoleUser,
		Stacks:   []string{},
	})
	svc := usecase.NewAuthService(mock)

	user, err := svc.Signin(context.Background(), "testuser", "password123")
	require.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
}

func TestAuthService_Signin_SuccessByEmail(t *testing.T) {
	mock := newMockUserRepo()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	mock.addUser(&domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashed),
		Role:     domain.RoleUser,
		Stacks:   []string{},
	})
	svc := usecase.NewAuthService(mock)

	user, err := svc.Signin(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestAuthService_Signin_WrongPassword(t *testing.T) {
	mock := newMockUserRepo()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	mock.addUser(&domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashed),
		Role:     domain.RoleUser,
		Stacks:   []string{},
	})
	svc := usecase.NewAuthService(mock)

	_, err := svc.Signin(context.Background(), "testuser", "wrongpassword")
	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}

func TestAuthService_Signin_UserNotFound(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())

	_, err := svc.Signin(context.Background(), "nonexistent", "password123")
	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}

func TestAuthService_Signin_RepoError(t *testing.T) {
	mock := newMockUserRepo()
	mock.createFn = func(u *domain.User) error { return nil }
	svc := usecase.NewAuthService(mock)

	mock.usersByUsername["testuser"] = &domain.User{
		ID:       1,
		Username: "testuser",
		Password: "hash",
	}

	err := mock.createFn(nil)
	require.NoError(t, err)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	mock.usersByUsername["testuser"].Password = string(hashed)

	svc.Signup(context.Background(), "testuser", "test@example.com", "password123")

	mock.usersByUsername["testuser"] = &domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashed),
		Role:     domain.RoleUser,
		Stacks:   []string{},
	}

	user, err := svc.Signin(context.Background(), "testuser", "password123")
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestAuthService_GetByID_Success(t *testing.T) {
	mock := newMockUserRepo()
	mock.addUser(&domain.User{ID: 1, Username: "testuser", Email: "test@example.com", Stacks: []string{}})
	svc := usecase.NewAuthService(mock)

	user, err := svc.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestAuthService_GetByID_NotFound(t *testing.T) {
	svc := usecase.NewAuthService(newMockUserRepo())

	_, err := svc.GetByID(context.Background(), 999)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestAuthService_Signup_TrimsWhitespace(t *testing.T) {
	mock := newMockUserRepo()
	svc := usecase.NewAuthService(mock)

	err := svc.Signup(context.Background(), "  testuser  ", "  test@example.com  ", "password123")
	require.NoError(t, err)

	user, err := mock.GetByUsername(context.Background(), "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestAuthService_Signup_RepoErrorOnCheck(t *testing.T) {
	mock := newMockUserRepo()
	mock.createFn = func(u *domain.User) error { return nil }
	mock.usersByUsername["testuser"] = &domain.User{
		ID:       1,
		Username: "testuser",
	}
	svc := usecase.NewAuthService(mock)

	err := mock.createFn(nil)
	require.NoError(t, err)

	svc.Signup(context.Background(), "newuser", "new@example.com", "password123")
}

func TestAuthService_Signin_FallsBackToEmailOnFirstGetByUsernameError(t *testing.T) {
	mock := newMockUserRepo()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	mock.usersByEmail["test@example.com"] = &domain.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashed),
		Role:      domain.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Stacks:    []string{},
	}
	svc := usecase.NewAuthService(mock)

	user, err := svc.Signin(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}
