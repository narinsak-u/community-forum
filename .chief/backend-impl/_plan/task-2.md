# Task Spec: task-2 — Session Auth
## Overview
Implement session-based authentication using Fiber sessions middleware with PostgreSQL as session store. Handlers: signup (bcrypt), signin, signout, me. Middleware: auth guard for protected routes.
## Implementation Steps
### 2.1 — Install Dependencies
`go get github.com/gofiber/session/v2` and `go get golang.org/x/crypto/bcrypt`
### 2.2 — Create `internal/middleware/middleware.go`
```go
func RequireAuth(c *fiber.Ctx) error {
    store := session.GetSession(c)
    userID := store.Get("user_id")
    if userID == nil {
        return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
    }
    return c.Next()
}
func GetUserID(c *fiber.Ctx) uint {
    store := session.GetSession(c)
    if id, ok := store.Get("user_id").(uint); ok { return id }
    return 0
}
func GetUserRole(c *fiber.Ctx) string {
    store := session.GetSession(c)
    if role, ok := store.Get("user_role").(string); ok { return role }
    return ""
}
```
### 2.3 — Create `internal/handlers/auth.go`
**SignupHandler**: parse `{ username, email, password }`, validate (username 3-30 chars alphanumeric, valid email, password min 8 chars), check duplicates (409), bcrypt hash (cost 12), create user (201).
**SigninHandler**: parse `{ login, password }`, find by username OR email, bcrypt compare (401 on failure), create session, set `session. Set("user_id", user.ID)`, set role in session, save, return user data (200).
**SignoutHandler**: get session, destroy, return 200.
**MeHandler**: get user ID from session, fetch from DB, return user data.
### 2.4 — Update `cmd/server/main.go`
1. Initialize Fiber sessions middleware: `session.New()` with config `{ KeyPrefix: "sess_" }`.
2. Add env var `SESSION_STORE` (default: `inmemory`, options: `inmemory`, `database`).
3. If `SESSION_STORE=database`, configure PostgreSQL store backed by the `sessions` table.
4. Public routes: `POST /api/v1/auth/signup`, `POST /api/v1/auth/signin`.
5. Protected routes with `RequireAuth`: `POST /api/v1/auth/signout`, `GET /api/v1/auth/me`, all thread/comment routes.
6. Keep existing stub routes during transition.
## Acceptance Criteria
1. `POST /api/v1/auth/signup` → 201 with bcrypt-hashed password, no password in response.
2. `POST /api/v1/auth/signin` (valid) → 200 + session cookie.
3. `POST /api/v1/auth/signin` (invalid) → 401.
4. `GET /api/v1/auth/me` (valid session) → user data.
5. `GET /api/v1/auth/me` (no session) → 401.
6. Duplicate username/email on signup → 409.
7. `go build ./cmd/server` succeeds.
8. Session cookie is HttpOnly.
## Edge Cases
- Login accepts both username and email: `db.Where("username = ? OR email = ?", login, login)`.
- Case-insensitive email: normalize to lowercase before storing.
- Bcrypt cost 12: add env var `BCRYPT_COST` (default 12) for CI environments where cost 12 is too slow.
- Session expiry: set `ExpiresAt` on session records. Add `// TODO: cleanup goroutine for expired sessions`.
- In-memory store for v1: fine for local dev. Document that PostgreSQL store is needed for production.