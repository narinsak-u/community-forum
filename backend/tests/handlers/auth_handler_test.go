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
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthApp(authService ports.AuthService) *fiber.App {
	app := fiber.New()
	authHandler := handlers.NewAuthHandler(authService)

	app.Post("/api/signup", authHandler.SignupHandler)
	app.Post("/api/signin", authHandler.SigninHandler)
	app.Post("/api/signout", authHandler.SignoutHandler)
	app.Get("/api/me", authHandler.MeHandler)

	return app
}

func TestSignupHandler_Success(t *testing.T) {
	svc := &mockAuthService{
		signupFn: func(ctx context.Context, username, email, password string) error {
			assert.Equal(t, "newuser", username)
			assert.Equal(t, "new@example.com", email)
			assert.Equal(t, "password123", password)
			return nil
		},
	}
	app := setupAuthApp(svc)

	body := `{"username":"newuser","email":"new@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "User registered successfully", result["message"])
}

func TestSignupHandler_InvalidBody(t *testing.T) {
	app := setupAuthApp(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/api/signup", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSignupHandler_ServiceError(t *testing.T) {
	svc := &mockAuthService{
		signupFn: func(ctx context.Context, username, email, password string) error {
			return fmt.Errorf("username already taken")
		},
	}
	app := setupAuthApp(svc)

	body := `{"username":"taken","email":"taken@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "username already taken", result["error"])
}

func TestSigninHandler_Success(t *testing.T) {
	middleware.InitSessionStore()

	svc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Username: "johndoe",
				Email:    "john@example.com",
				Avatar:   "avatar.jpg",
				Bio:      "Hello!",
				Role:     domain.RoleUser,
			}, nil
		},
	}
	app := setupAuthApp(svc)

	body := `{"login":"johndoe","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "johndoe", result["username"])
	assert.Equal(t, float64(1), result["id"])
}

func TestSigninHandler_InvalidCredentials(t *testing.T) {
	middleware.InitSessionStore()

	svc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return nil, fmt.Errorf("invalid credentials")
		},
	}
	app := setupAuthApp(svc)

	body := `{"login":"wrong","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSigninHandler_InvalidBody(t *testing.T) {
	middleware.InitSessionStore()
	app := setupAuthApp(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`bad json`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSignoutHandler_Success(t *testing.T) {
	middleware.InitSessionStore()
	app := setupAuthApp(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/api/signout", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMeHandler_Unauthorized(t *testing.T) {
	middleware.InitSessionStore()
	app := setupAuthApp(&mockAuthService{})

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestMeHandler_Authorized(t *testing.T) {
	middleware.InitSessionStore()

	svc := &mockAuthService{
		signinFn: func(ctx context.Context, login, password string) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Username: "johndoe",
				Email:    "john@example.com",
				Role:     domain.RoleUser,
			}, nil
		},
		getByIDFn: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				Username: "johndoe",
				Email:    "john@example.com",
				Avatar:   "avatar.jpg",
				Bio:      "Hello!",
				Role:     domain.RoleUser,
			}, nil
		},
	}
	app := setupAuthApp(svc)

	signinBody := `{"login":"johndoe","password":"password123"}`
	signinReq := httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(signinBody))
	signinReq.Header.Set("Content-Type", "application/json")
	signinResp, err := app.Test(signinReq)
	require.NoError(t, err)

	cookies := signinResp.Header["Set-Cookie"]
	require.NotEmpty(t, cookies)

	meReq := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	meReq.Header.Set("Cookie", strings.Join(cookies, "; "))

	meResp, err := app.Test(meReq)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, meResp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(meResp.Body).Decode(&result)
	assert.Equal(t, "johndoe", result["username"])
}
