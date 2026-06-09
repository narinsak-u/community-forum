package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db_adapter "community-forum/backend/internal/adapters/db"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCommentApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	sessionManager := setupTestSessionManager()
	commentRepo := db_adapter.NewGORMCommentRepository(db)
	threadRepo := db_adapter.NewGORMThreadRepository(db)
	commentService := usecase.NewCommentService(commentRepo, threadRepo)
	commentHandler := handlers.NewCommentHandler(commentService, sessionManager)

	app.Post("/api/threads/:slug/comments", commentHandler.CreateCommentHandler)
	app.Delete("/api/comments/:id", commentHandler.DeleteCommentHandler)

	return app
}

func TestCreateCommentHandler_InvalidBody(t *testing.T) {
	app := setupCommentApp(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/comments", strings.NewReader(`bad json`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateCommentHandler_ShortContent(t *testing.T) {
	app := setupCommentApp(nil)

	body := `{"content":"x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateCommentHandler_LongContent(t *testing.T) {
	app := setupCommentApp(nil)

	longContent := strings.Repeat("a", 10001)
	body := `{"content":"` + longContent + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteCommentHandler_InvalidID(t *testing.T) {
	app := setupCommentApp(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/comments/abc", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
