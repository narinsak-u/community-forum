package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// SessionManager handles creating, retrieving, and saving session data.
type SessionManager struct {
	Store *session.Store
}

// NewSessionManager initializes the session manager configuration.
func NewSessionManager() *SessionManager {
	store := session.New(session.Config{
		CookieHTTPOnly: true,  // Prevents JavaScript from accessing the cookie (security against XSS)
		CookieSameSite: "Lax", // Helps prevent CSRF attacks
		CookieSecure:   false, // Set to true in production to only allow sessions over HTTPS
	})
	return &SessionManager{Store: store}
}

// RequireAuth is a middleware function that stops a request if the user is not logged in.
func (m *SessionManager) RequireAuth(c *fiber.Ctx) error {
	// Try to get the session from the request context (via cookies)
	sess, err := m.Store.Get(c)

	// If the session doesn't exist or doesn't have a user_id, return 401 Unauthorized.
	if err != nil || sess.Get("user_id") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// c.Next() tells Fiber to move to the next function in the chain (the actual handler).
	return c.Next()
}

// GetUserID is a helper to retrieve the logged-in user's ID from the session.
func (m *SessionManager) GetUserID(c *fiber.Ctx) uint {
	sess, err := m.Store.Get(c)
	if err != nil {
		return 0
	}

	// Type assertion: session data is stored as interface{}, we need to cast it back to uint.
	if id, ok := sess.Get("user_id").(uint); ok {
		return id
	}
	return 0
}

// GetUserRole is a helper to retrieve the logged-in user's role (e.g., 'admin', 'user').
func (m *SessionManager) GetUserRole(c *fiber.Ctx) string {
	sess, err := m.Store.Get(c)
	if err != nil {
		return ""
	}

	if role, ok := sess.Get("user_role").(string); ok {
		return role
	}
	return ""
}
