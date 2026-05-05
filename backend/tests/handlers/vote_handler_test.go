package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupVoteApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	voteHandler := handlers.NewVoteHandler(db)

	app.Post("/api/threads/:slug/vote", voteHandler.VoteThreadHandler)
	app.Post("/api/comments/:id/vote", voteHandler.VoteCommentHandler)

	return app
}

func TestVoteThreadHandler_InvalidBody(t *testing.T) {
	middleware.InitSessionStore()
	app := setupVoteApp(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/vote", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteThreadHandler_InvalidValue(t *testing.T) {
	middleware.InitSessionStore()
	app := setupVoteApp(nil)

	body := `{"value":99}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/test-thread/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidBody(t *testing.T) {
	middleware.InitSessionStore()
	app := setupVoteApp(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/comments/1/vote", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidCommentID(t *testing.T) {
	middleware.InitSessionStore()
	app := setupVoteApp(nil)

	body := `{"value":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/comments/abc/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestVoteCommentHandler_InvalidValue(t *testing.T) {
	middleware.InitSessionStore()
	app := setupVoteApp(nil)

	body := `{"value":99}`
	req := httptest.NewRequest(http.MethodPost, "/api/comments/1/vote", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
