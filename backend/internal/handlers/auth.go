package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler manages authentication-related requests like login, signup, and logout.
type AuthHandler struct {
	// AuthService is an interface (a port in Hexagonal Architecture).
	// This allows the handler to call business logic without knowing the implementation details.
	AuthService ports.AuthService
}

// NewAuthHandler is a constructor that "injects" the necessary services into the handler.
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Request structs define the expected JSON structure from the client.
// The `json:"..."` tags tell Fiber how to map JSON keys to Go fields.

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
	
	// c.BodyParser automatically reads the JSON request and fills our 'req' struct.
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Call the service layer to perform the actual signup logic.
	if err := h.AuthService.Signup(c.Context(), req.Username, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return 201 Created on success.
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

	// Signin returns the user domain object if credentials are correct.
	user, err := h.AuthService.Signin(c.Context(), req.Login, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Retrieve the user's session from the session store.
	sess, err := middleware.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session",
		})
	}

	// Store user information in the session so we know who they are in future requests.
	sess.Set("user_id", user.ID)
	sess.Set("user_role", string(user.Role))

	// Save the session changes (this typically sends a cookie back to the user).
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	// Send back the user details (omitting sensitive info like password).
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
	sess, err := middleware.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get session",
		})
	}

	// Destroy removes the session from the store and clears the client's cookie.
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
	// middleware.GetUserID is a helper that reads the user_id from the session.
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Fetch the most up-to-date user data from the database.
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
