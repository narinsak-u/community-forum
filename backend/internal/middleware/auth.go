package middleware

import (
	"strings"
	"time"

	"community-forum/backend/internal/lib"

	"github.com/gofiber/fiber/v2"
)

const jwtCookieName = "forge_token"

type SessionManager struct {
	Secret string
	Expiry time.Duration
}

func NewSessionManager(secret string, expiry time.Duration) *SessionManager {
	return &SessionManager{Secret: secret, Expiry: expiry}
}

func (m *SessionManager) SignToken(userID uint, role string) (string, error) {
	return lib.SignJWT(userID, role, m.Secret, m.Expiry)
}

func (m *SessionManager) SetTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
	})
}

func (m *SessionManager) ClearTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",
		MaxAge:   -1,
	})
}

func (m *SessionManager) RequireAuth(c *fiber.Ctx) error {
	tokenStr := m.extractToken(c)
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, err := lib.VerifyJWT(tokenStr, m.Secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("user_role", claims.Role)
	return c.Next()
}

func (m *SessionManager) GetUserID(c *fiber.Ctx) uint {
	if id, ok := c.Locals("user_id").(uint); ok {
		return id
	}

	// Fallback: parse from JWT cookie even if RequireAuth wasn't called
	tokenStr := m.extractToken(c)
	if tokenStr == "" {
		return 0
	}
	claims, err := lib.VerifyJWT(tokenStr, m.Secret)
	if err != nil {
		return 0
	}
	return claims.UserID
}

func (m *SessionManager) GetUserRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("user_role").(string); ok {
		return role
	}

	tokenStr := m.extractToken(c)
	if tokenStr == "" {
		return ""
	}
	claims, err := lib.VerifyJWT(tokenStr, m.Secret)
	if err != nil {
		return ""
	}
	return claims.Role
}

func (m *SessionManager) extractToken(c *fiber.Ctx) string {
	cookie := c.Cookies(jwtCookieName)
	if cookie != "" {
		return cookie
	}

	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}
