package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Store is the global session manager. 
// It handles creating, retrieving, and saving session data (like cookies).
var Store *session.Store

// InitSessionStore sets up the session configuration.
func InitSessionStore() {
	Store = session.New(session.Config{
		CookieHTTPOnly: true,   // Prevents JavaScript from accessing the cookie (security against XSS)
		CookieSameSite: "Lax",    // Helps prevent CSRF attacks
		CookieSecure:   false,  // Set to true in production to only allow sessions over HTTPS
	})
}

// RequireAuth is a middleware function that stops a request if the user is not logged in.
// In Fiber, middleware is a function that runs before the actual handler.
func RequireAuth(c *fiber.Ctx) error {
	// Try to get the session from the request context (via cookies)
	sess, err := Store.Get(c)
	
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
func GetUserID(c *fiber.Ctx) uint {
	sess, err := Store.Get(c)
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
func GetUserRole(c *fiber.Ctx) string {
	sess, err := Store.Get(c)
	if err != nil {
		return ""
	}
	
	if role, ok := sess.Get("user_role").(string); ok {
		return role
	}
	return ""
}
