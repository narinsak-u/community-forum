package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func InitSessionStore() {
	Store = session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		CookieSecure:   false,
	})
}

func RequireAuth(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil || sess.Get("user_id") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return c.Next()
}

func GetUserID(c *fiber.Ctx) uint {
	sess, err := Store.Get(c)
	if err != nil {
		return 0
	}
	if id, ok := sess.Get("user_id").(uint); ok {
		return id
	}
	return 0
}

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
