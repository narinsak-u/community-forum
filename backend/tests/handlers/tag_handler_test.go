package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTagApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	tagHandler := handlers.NewTagHandler(db)

	app.Get("/api/tags", tagHandler.ListTagsHandler)
	app.Post("/api/tags", tagHandler.CreateTagHandler)

	return app
}

func TestListTagsHandler_Success(t *testing.T) {
	gormDB, mock := newMockDB(t)
	app := setupTagApp(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE "tags"."deleted_at" IS NULL ORDER BY name ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "color"}).AddRow(1, "Go", "#00ADD8").AddRow(2, "React", "#61DAFB"))

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	tags := result["tags"].([]interface{})
	assert.Len(t, tags, 2)
}

func TestCreateTagHandler_Unauthorized(t *testing.T) {
	middleware.InitSessionStore()
	app := setupTagApp(nil)

	body := `{"name":"NewTag","color":"#ff0000"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestCreateTagHandler_AdminSuccess(t *testing.T) {
	middleware.InitSessionStore()
	gormDB, mock := newMockDB(t)
	app := setupTagApp(gormDB)

	signinSvc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: "admin", Email: "admin@test.com", Role: domain.RoleAdmin}, nil
		},
	}
	authApp := setupAuthApp(signinSvc)
	signinReq := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"login":"admin","password":"pass"}`))
	signinReq.Header.Set("Content-Type", "application/json")
	signinResp, _ := authApp.Test(signinReq)
	cookies := signinResp.Header["Set-Cookie"]
	if len(cookies) == 0 {
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags"`)).
		WithArgs("golang").
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tags"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	createReq := httptest.NewRequest(http.MethodPost, "/api/tags", strings.NewReader(`{"name":"golang","color":"#00FF00"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Cookie", strings.Join(cookies, "; "))

	resp, err := app.Test(createReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateTagHandler_InvalidBody(t *testing.T) {
	middleware.InitSessionStore()
	gormDB, _ := newMockDB(t)
	app := setupTagApp(gormDB)

	signinSvc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: "admin", Email: "admin@test.com", Role: domain.RoleAdmin}, nil
		},
	}
	authApp := setupAuthApp(signinSvc)
	signinReq := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"login":"admin","password":"pass"}`))
	signinReq.Header.Set("Content-Type", "application/json")
	signinResp, _ := authApp.Test(signinReq)
	cookies := signinResp.Header["Set-Cookie"]
	if len(cookies) == 0 {
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/api/tags", strings.NewReader(`bad`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", strings.Join(cookies, "; "))

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateTagHandler_ColorValidation(t *testing.T) {
	middleware.InitSessionStore()
	gormDB, _ := newMockDB(t)
	app := setupTagApp(gormDB)

	signinSvc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{ID: 1, Username: "admin", Email: "admin@test.com", Role: domain.RoleAdmin}, nil
		},
	}
	authApp := setupAuthApp(signinSvc)
	signinReq := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"login":"admin","password":"pass"}`))
	signinReq.Header.Set("Content-Type", "application/json")
	signinResp, _ := authApp.Test(signinReq)
	cookies := signinResp.Header["Set-Cookie"]
	if len(cookies) == 0 {
		return
	}

	body := `{"name":"MyTag","color":"invalid-color"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", strings.Join(cookies, "; "))

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
