package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserApp(userService *mockUserService, db *gorm.DB) *fiber.App {
	app := fiber.New()
	userHandler := handlers.NewUserHandler(userService, db)
	userHandler.UserService = userService
	userHandler.DB = db

	app.Get("/api/users/:username", userHandler.GetUserHandler)
	app.Put("/api/users/:username", userHandler.UpdateUserHandler)
	app.Get("/api/users/:username/threads", userHandler.GetUserThreadsHandler)

	return app
}

func TestGetUserHandler_Success(t *testing.T) {
	svc := &mockUserService{
		getUserProfileFn: func(ctx context.Context, username string) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Username: "johndoe",
				Avatar:   "avatar.jpg",
				Bio:      "Hello!",
				Stacks:   []string{"Go", "React"},
				Role:     domain.RoleUser,
			}, nil
		},
	}
	app := setupUserApp(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/users/johndoe", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	user := result["user"].(map[string]interface{})
	assert.Equal(t, "johndoe", user["username"])
}

func TestGetUserHandler_NotFound(t *testing.T) {
	svc := &mockUserService{
		getUserProfileFn: func(ctx context.Context, username string) (*domain.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}
	app := setupUserApp(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/users/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateUserHandler_Unauthorized(t *testing.T) {
	middleware.InitSessionStore()
	app := setupUserApp(&mockUserService{}, nil)

	body := `{"bio":"new bio"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/johndoe", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUpdateUserHandler_Forbidden(t *testing.T) {
	middleware.InitSessionStore()
	otherSvc := &mockUserService{
		getUserProfileFn: func(ctx context.Context, username string) (*domain.User, error) {
			return &domain.User{ID: 2, Username: "other"}, nil
		},
	}
	app := setupUserApp(otherSvc, nil)

	signinBody := `{"login":"johndoe","password":"pass"}`
	signinReq := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(signinBody))
	signinReq.Header.Set("Content-Type", "application/json")
	signinResp, _ := app.Test(signinReq)

	svc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: "johndoe", Email: "john@test.com", Role: domain.RoleUser}, nil
		},
	}
	authApp := setupAuthApp(svc)
	signinReq2 := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(signinBody))
	signinReq2.Header.Set("Content-Type", "application/json")
	signinResp, _ = authApp.Test(signinReq2)

	cookies := signinResp.Header["Set-Cookie"]
	if len(cookies) == 0 {
		return
	}

	otherSvc.getUserProfileFn = func(ctx context.Context, username string) (*domain.User, error) {
		return &domain.User{ID: 2, Username: "otheruser"}, nil
	}

	updateReq := httptest.NewRequest(http.MethodPut, "/api/users/otheruser", strings.NewReader(`{"bio":"hacked"}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Cookie", strings.Join(cookies, "; "))

	resp, err := app.Test(updateReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestGetUserThreadsHandler_Success(t *testing.T) {
	middleware.InitSessionStore()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`).
		WithArgs("johndoe").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "johndoe"))

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE author_id = \$1 AND "threads"\."deleted_at" IS NULL`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	app := setupUserApp(&mockUserService{}, gormDB)

	req := httptest.NewRequest(http.MethodGet, "/api/users/johndoe/threads?page=1&pageSize=5", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetUserThreadsHandler_UserNotFound(t *testing.T) {
	middleware.InitSessionStore()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`).
		WithArgs("nonexistent").
		WillReturnError(gorm.ErrRecordNotFound)

	app := setupUserApp(&mockUserService{}, gormDB)

	req := httptest.NewRequest(http.MethodGet, "/api/users/nonexistent/threads", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetUserThreadsHandler_PaginationDefaults(t *testing.T) {
	middleware.InitSessionStore()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`).
		WithArgs("johndoe").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "johndoe"))

	mock.ExpectQuery(`SELECT count\(\*\) FROM "threads" WHERE author_id = \$1 AND "threads"\."deleted_at" IS NULL`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	app := setupUserApp(&mockUserService{}, gormDB)

	req := httptest.NewRequest(http.MethodGet, "/api/users/johndoe/threads", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
