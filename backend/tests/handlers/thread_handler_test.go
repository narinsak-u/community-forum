package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	db_adapter "community-forum/backend/internal/adapters/db"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/usecase"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupThreadApp(db *gorm.DB) (*fiber.App, *middleware.SessionManager) {
	app := fiber.New()
	sessionManager := setupTestSessionManager()
	threadRepo := db_adapter.NewGORMThreadRepository(db)
	threadService := usecase.NewThreadService(threadRepo)
	threadHandler := handlers.NewThreadHandler(threadService, sessionManager)

	app.Post("/api/threads", threadHandler.CreateThreadHandler)
	app.Get("/api/threads", threadHandler.ListThreadsHandler)
	app.Get("/api/threads/featured", threadHandler.FeaturedThreadHandler)
	app.Get("/api/threads/trending", threadHandler.TrendingThreadsHandler)
	app.Get("/api/threads/:slug", threadHandler.GetThreadHandler)
	app.Put("/api/threads/:slug", threadHandler.UpdateThreadHandler)
	app.Delete("/api/threads/:slug", threadHandler.DeleteThreadHandler)

	return app, sessionManager
}

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)
	return gormDB, mock
}

func TestCreateThreadHandler_InvalidTitle(t *testing.T) {
	gormDB, _ := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	tests := []struct {
		name string
		body string
	}{
		{"too short", `{"title":"abc","content":"this is content that is long enough","status":"published"}`},
		{"too long", `{"title":"` + strings.Repeat("a", 256) + `","content":"long content here for the test","status":"published"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestCreateThreadHandler_InvalidBody(t *testing.T) {
	gormDB, _ := newMockDB(t)
	app, _ := setupThreadApp(gormDB)
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestListThreadsHandler_Pagination(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "threads"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "status", "author_id", "view_count", "created_at", "upvotes", "downvotes", "replies_count"}))

	req := httptest.NewRequest(http.MethodGet, "/api/threads?page=2&pageSize=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	pagination := result["pagination"].(map[string]interface{})
	assert.Equal(t, float64(2), pagination["page"])
	assert.Equal(t, float64(10), pagination["pageSize"])
	assert.Equal(t, float64(0), pagination["total"])
	assert.Equal(t, float64(0), pagination["totalPages"])
}

func TestGetThreadHandler_Success(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WithArgs("hello-world").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "status", "author_id"}).
			AddRow(1, "Hello World", "hello-world", "Content here", "published", 1))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "threads" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "status", "author_id", "view_count"}).
			AddRow(1, "Hello World", "hello-world", "Content here", "published", 1, 5))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "avatar"}).AddRow(1, "johndoe", ""))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "thread_tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"thread_id", "tag_id"}))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments"`)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	req := httptest.NewRequest(http.MethodGet, "/api/threads/hello-world", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetThreadHandler_NotFound(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WithArgs("nonexistent").
		WillReturnError(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestFeaturedThreadHandler_NotFound(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/featured", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestTrendingThreadsHandler_DBError(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/threads/trending", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteThreadHandler_NotFound(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WithArgs("nonexistent").
		WillReturnError(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/threads/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListThreadsHandler_SortOldest(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "threads"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`SELECT .+ FROM "threads"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "slug", "content", "status", "author_id", "view_count", "created_at", "upvotes", "downvotes", "replies_count"}))

	req := httptest.NewRequest(http.MethodGet, "/api/threads?sort=oldest&page=1&pageSize=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateThreadHandler_InvalidStatus(t *testing.T) {
	gormDB, _ := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	body := `{"title":"Valid Title Here","content":"This is valid content for the thread body","status":"invalid_status"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateThreadHandler_ShortContent(t *testing.T) {
	gormDB, _ := newMockDB(t)
	app, _ := setupThreadApp(gormDB)

	body := `{"title":"Valid Title","content":"short","status":"published"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
