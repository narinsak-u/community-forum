package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler manages authentication-related requests like login, signup, and logout.
type AuthHandler struct {
	AuthService    ports.AuthService
	SessionManager *middleware.SessionManager
}

// NewAuthHandler is a constructor that "injects" the necessary services into the handler.
func NewAuthHandler(authService ports.AuthService, sessionManager *middleware.SessionManager) *AuthHandler {
	return &AuthHandler{
		AuthService:    authService,
		SessionManager: sessionManager,
	}
}

// Request structs define the expected JSON structure from the client.
type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// SignupHandler handles POST requests to create a new user account.
func (h *AuthHandler) SignupHandler(c *fiber.Ctx) error {
	var req SignupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.AuthService.Signup(c.Context(), req.Username, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

// SigninHandler handles POST requests to log a user in.
func (h *AuthHandler) SigninHandler(c *fiber.Ctx) error {
	var req SigninRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.AuthService.Signin(c.Context(), req.Login, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	sess, err := h.SessionManager.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_role", string(user.Role))

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"role":       user.Role,
		"stacks":     user.Stacks,
		"created_at": user.CreatedAt,
	})
}

// SignoutHandler handles POST requests to log the user out.
func (h *AuthHandler) SignoutHandler(c *fiber.Ctx) error {
	sess, err := h.SessionManager.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get session",
		})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to destroy session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// MeHandler returns the currently logged-in user's profile.
func (h *AuthHandler) MeHandler(c *fiber.Ctx) error {
	userID := h.SessionManager.GetUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, err := h.AuthService.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"role":       user.Role,
		"stacks":     user.Stacks,
		"created_at": user.CreatedAt,
	})
}
