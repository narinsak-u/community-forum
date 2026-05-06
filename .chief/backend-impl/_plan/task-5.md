# Task Spec: task-5 — User, Tag Handlers & Main Refactor
## Overview
Implement user profile CRUD, threads-by-user listing, tag CRUD. Refactor main.go to remove all inline handler logic, delegating everything to the handler service layer.
## Implementation Steps
### 5.1 — Create `internal/handlers/users.go`
**GetUserHandler** (public):
1. Lookup by username. Return `{ id, username, avatar, bio, stacks, role, created_at }` — no email or password.
2. If not found → 404.
**UpdateUserHandler** (protected):
1. Lookup by username. Verify ownership. If not → 403.
2. Parse body: `{ avatar, bio, stacks }`. Partial update (only non-nil fields).
3. If `stacks` set: JSON-encode string array to `User.Stacks`.
4. Save, return updated user.
**GetUserThreadsHandler**:
1. Lookup user by username. If not found → 404.
2. Query: threads by `author_id = user.ID`, ordered by `created_at DESC`.
3. Pagination (page, pageSize, sort). Preload `Author`, `Tags`.
4. Compute vote counts. Return threads + pagination.
### 5.2 — Create `internal/handlers/tags.go`
**ListTagsHandler** (public):
1. Fetch all tags ordered by name.
2. Return `{ tags: [...] }`.
**CreateTagHandler** (protected + admin):
1. Check admin role: fetch user from DB, verify `Role == "admin"`. If not → 403.
2. Parse body: `{ name, color }`. Validate: name 3-50 chars, color valid hex (default `#6b366f1`).
3. Check duplicate name (case-insensitive). If exists → 409.
4. Create tag. Return 201.
### 5.3 — Main Refactor
In `cmd/server/main.go`:
1. Remove ALL inline `func(c *fiber.Ctx) error` blocks.
2. Keep: env loading, DB connect, AutoMigrate, middleware (recover, logger, cors, sessions), route registrations using handler functions.
3. Keep `/` and `/health` routes.
4. Run `go fmt ./... && go vet ./...` — zero errors required.
### 5.4 — Write Tests
Write at least one test per handler (basic happy path):
- `internal/handlers/auth_test.go`: signup, signin
- `internal/handlers/threads_test.go`: create, list
- `internal/handlers/votes_test.go`: vote
### 5.5 — Vote Helper Refinement
Update all handler references to vote helpers. Pass `db` from the handler (create with `db. Model(&Thread{})` if needed).
## Acceptance Criteria
1. Get user profile → 200, no email or password in response.
2. Update own profile → 200.
3. Update another's profile → 403.
4. Get user threads → paginated list.
5. List tags → all tags.
6. Create tag (admin) → 201.
7. Create tag (non-admin) → 403.
8. main.go: zero inline handlers.
9. `go fmt ./... && go vet ./...` → zero errors.
10. `go test ./...` → all tests pass.
## Edge Cases
- Stacks encoding: store as JSON in `text`/`jsonb` column. Test with `["Go", "TypeScript", "Rust"]`.
- Admin role check: fetch user from DB (`GetUserRole` helper uses session, not DB). Use `db.First(&user, id)` or cache in session.
- Tag name case-insensitivity: use `db.Where("LOWER(name) = LOWER(?)", name)`.
- Username uniqueness: existing unique index sufficient.