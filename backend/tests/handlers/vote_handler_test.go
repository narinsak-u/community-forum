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

func setupVoteApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	sessionManager := setupTestSessionManager()
	voteRepo := db_adapter.NewGORMVoteRepository(db)
	threadRepo := db_adapter.NewGORMThreadRepository(db)
	commentRepo := db_adapter.NewGORMCommentRepository(db)
	voteService := usecase.NewVoteService(voteRepo, threadRepo, commentRepo)
	voteHandler := handlers.NewVoteHandler(voteService, sessionManager)

	app.Post("/api/threads/:slug/vote", voteHandler.VoteThreadHandler)
	app.Post("/api/comments/:id/vote", voteHandler.VoteCommentHandler)

	return app
}

func TestVoteThreadHandler_InvalidBody(t *testing.T) {
	app := setupVoteApp(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/vote", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteThreadHandler_InvalidValue(t *testing.T) {
	app := setupVoteApp(nil)

	body := `{"value":99}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidBody(t *testing.T) {
	app := setupVoteApp(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/comments/1/vote", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidCommentID(t *testing.T) {
	app := setupVoteApp(nil)

	body := `{"value":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/comments/abc/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidValue(t *testing.T) {
	app := setupVoteApp(nil)

	body := `{"value":99}`
	req := httptest.NewRequest(http.MethodPost, "/api/comments/1/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
