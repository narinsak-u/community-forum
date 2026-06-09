# Backend Workflows

> Detailed step-by-step documentation of the runtime behavior of the Community
> Forum backend. Last updated: 2026-06-08.
>
> All file references use `path:line` notation relative to the `backend/`
> directory. See [`OVERVIEW.md`](./OVERVIEW.md) for the structural overview.

---

## Table of Contents

1. [Server startup flow](#1-server-startup-flow)
2. [Database connection and migration flow](#2-database-connection-and-migration-flow)
3. [Authentication flow](#3-authentication-flow)
   - [Registration](#31-registration)
   - [Login](#32-login)
   - [Token validation on protected routes](#33-token-validation-on-protected-routes)
   - [Logout](#34-logout)
4. [Request lifecycle](#4-request-lifecycle)
5. [CRUD workflows](#5-crud-workflows)
   - [Threads (create / read / update / delete)](#51-threads)
   - [Comments (create / read-nested / delete)](#52-comments)
   - [Votes (upvote / downvote / retract)](#53-votes)
   - [User profile (read / update)](#54-user-profile)
   - [Tags (list / create)](#55-tags)
6. [Error handling flow](#6-error-handling-flow)
7. [Cross-cutting concerns](#7-cross-cutting-concerns)

---

## 1. Server startup flow

```
main() (cmd/server/main.go:29)
   |
   +-- 1. cfg := config.Load()                         (main.go:30)
   |      +-- internal/config/config.go:30-47
   |         - godotenv.Load() -- read .env (or warn and continue)
   |         - Build *Config with env vars + defaults
   |
   +-- 2. db := config.InitDB(cfg)                     (main.go:31)
   |      +-- see Database connection flow (section 2)
   |
   +-- 3. app := fiber.New(fiber.Config{ ErrorHandler })  (main.go:32-40)
   |      - Custom ErrorHandler returns JSON { "error": "..." } with 500
   |        for any error returned from a handler
   |
   +-- 4. sessionManager := middleware.NewSessionManager(...)  (main.go:43)
   |      - Wraps the JWT secret + 72h expiry
   |
   +-- 5. app.Use(recover.New())                       (main.go:49)
   +-- 5. app.Use(logger.New())                        (main.go:50)
   +-- 5. app.Use(cors.New(...))                       (main.go:51-56)
   |      - Origins: cfg.CORSOrigins
   |      - AllowCredentials: true
   |      - Methods: GET,POST,PUT,DELETE,PATCH,OPTIONS
   |      - Headers: Origin,Content-Type,Accept,Authorization
   |
   +-- 6. app.Get("/",  inline health)                 (main.go:59-61)
   +-- 6. app.Get("/health", inline health)            (main.go:63-65)
   |
   +-- 7. api := app.Group("/api/v1")                  (main.go:68)
   |
   +-- 8. Wire adapters (outbound):                    (main.go:72-76)
   |      - userRepo    := db.NewGORMUserRepository(db)
   |      - threadRepo  := db.NewGORMThreadRepository(db)
   |      - commentRepo := db.NewGORMCommentRepository(db)
   |      - voteRepo    := db.NewGORMVoteRepository(db)
   |      - tagRepo     := db.NewGORMTagRepository(db)
   |
   +-- 9. Wire use cases (services):                   (main.go:79-84)
   |      - authService    := usecase.NewAuthService(userRepo)
   |      - userService    := usecase.NewUserService(userRepo)
   |      - threadService  := usecase.NewThreadService(threadRepo)
   |      - commentService := usecase.NewCommentService(commentRepo, threadRepo)
   |      - voteService    := usecase.NewVoteService(voteRepo, threadRepo, commentRepo)
   |      - tagService     := usecase.NewTagService(tagRepo)
   |
   +-- 10. Wire handlers (inbound adapters):          (main.go:87-92)
   |      - authHandler    := handlers.NewAuthHandler(authService, sessionManager)
   |      - threadHandler  := handlers.NewThreadHandler(threadService, sessionManager)
   |      - commentHandler := handlers.NewCommentHandler(commentService, sessionManager)
   |      - voteHandler    := handlers.NewVoteHandler(voteService, sessionManager)
   |      - userHandler    := handlers.NewUserHandler(userService, threadService, sessionManager)
   |      - tagHandler     := handlers.NewTagHandler(tagService, sessionManager)
   |
   +-- 11. Register auth-rate limiter                  (main.go:96-103)
   |      - 10 req/min per IP, SkipSuccessfulRequests=true
   |
   +-- 12. Register all /api/v1/* routes               (main.go:104-129)
   |
   +-- 13. app.Listen(":" + cfg.Port)                  (main.go:133)
          - log.Printf("Server starting on http://localhost:%s", cfg.Port)
          - log.Fatalf on bind error
```

### Key invariant
Every `Handler` depends only on **ports** (`*Service` interfaces) and the
`SessionManager`; the wiring in `main.go` is the **only** place that knows
about concrete GORM repositories. This is the Hexagonal composition root.

---

## 2. Database connection and migration flow

```
config.InitDB(cfg)  (internal/config/config.go:76-81)
   |
   +-- 1. db := ConnectDB(cfg)                        (config.go:56-62)
   |      +-- Builds DSN via cfg.DSN()  (config.go:49-54)
   |      |     fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", ...)
   |      +-- gorm.Open(postgres.Open(dsn), &gorm.Config{})
   |            - On error: log.Fatalf("Failed to connect to database: %v", err)
   |            - Returns *gorm.DB
   |
   +-- 2. MigrateDB(db)                                (config.go:64-74)
   |      +-- db.AutoMigrate(
   |            &models.User{},
   |            &models.Thread{},
   |            &models.Comment{},
   |            &models.Tag{},
   |            &models.Vote{},
   |         )
   |            - GORM creates tables + indexes if missing
   |            - Does NOT alter existing columns (additive only)
   |            - On error: log.Fatalf
   |
   +-- 3. seed.Seed(db)                                (internal/seed/seed.go)
          +-- db.Model(&models.User{}).Count(&userCount)
          +-- If userCount > 0: return (idempotent)
          |
          +-- Bcrypt-hash "password123" (MinCost for speed) for all 4 users
          +-- Insert 4 users: architect (admin), cypher, nexus, oracle
          +-- Insert 6 tags: architecture, security, distributed-systems,
          |                cryptography, protocol-design, systems
          +-- Insert 10 threads (with realistic markdown content)
          |     and m2m-associate tags via db.Model(&thread).Association("Tags").Append(&tag)
          +-- Insert 12 comments
          +-- Insert 11 votes (mix of +1 / -1)
   |
   +-- 4. return db
```

### Notes
- GORM `AutoMigrate` is **non-destructive**: dropping a field from a model
  leaves the column in place. To remove columns you must drop them
  manually or add a versioned migration.
- The `User.Stacks` field is stored as a JSONB column (model tag
  `gorm:"type:jsonb"`) and is marshaled/unmarshaled in the repository
  (`adapters/db/user_repository.go:107-148`).
- The `Vote` table has two composite unique indexes (`idx_vote_thread` and
  `idx_vote_comment`) on `(user_id, thread_id)` and `(user_id, comment_id)`
  respectively. They are nullable on the other side so that PostgreSQL
  permits one vote per (user, target).

---


## 3. Authentication flow

### 3.1 Registration

```
HTTP POST /api/v1/auth/signup
   | (rate-limited: 10/min/IP, but SkipSuccessfulRequests=true)
   |
   v
authHandler.SignupHandler  (handlers/auth.go:33-51)
   |
   +-- 1. c.BodyParser(&req)  -- decodes {username, email, password}
   |      on parse error -> 400 {"error":"Invalid request body"}
   |
   +-- 2. h.AuthService.Signup(ctx, username, email, password)
   |      (usecase/auth_service.go:45-94)
   |      |
   |      +-- Trim whitespace from username and email
   |      +-- Validate:
   |      |     - username 3-30 chars, regex ^[a-zA-Z0-9_]+$
   |      |     - email matches emailRegex
   |      |     - password >= 8 chars
   |      |     Any failure -> return plain error (handler maps to 400)
   |      |
   |      +-- repo.GetByUsername(ctx, username)
   |      |     - If found -> return ErrUsernameTaken
   |      +-- repo.GetByEmail(ctx, email)
   |      |     - If found -> return ErrEmailRegistered
   |      |
   |      +-- bcrypt.GenerateFromPassword([]byte(password), 12)
   |      |     - cost factor 12 (~250ms) -- strong default
   |      |
   |      +-- repo.Create(ctx, &domain.User{ ..., Role: domain.RoleUser })
   |            (adapters/db/user_repository.go:29-38)
   |            - toModel(u) JSON-marshals Stacks ([]string -> string)
   |            - r.db.WithContext(ctx).Create(m).Error
   |
   +-- 3. Return 201 {"message":"User registered successfully"}
```

### 3.2 Login

```
HTTP POST /api/v1/auth/signin  (with JSON { login, password })
   | (rate-limited)
   |
   v
authHandler.SigninHandler  (handlers/auth.go:53-87)
   |
   +-- 1. c.BodyParser(&req)  -> 400 on error
   |
   +-- 2. h.AuthService.Signin(ctx, login, password)
   |      (usecase/auth_service.go:97-123)
   |      +-- repo.GetByUsername(ctx, login)
   |      |     - if ErrNotFound -> try repo.GetByEmail(ctx, login)
   |      |     - if both fail -> return ErrInvalidCredentials
   |      +-- bcrypt.CompareHashAndPassword(hash, password)
   |      |     - mismatch -> ErrInvalidCredentials
   |      +-- return *domain.User
   |
   +-- 3. token, err := h.SessionManager.SignToken(user.ID, string(user.Role))
   |      (middleware/auth.go:23-25 -> lib.SignJWT)
   |      - Claims: { user_id, role, exp=now+72h, iat=now }
   |      - HS256 with cfg.JWTSecret
   |
   +-- 4. h.SessionManager.SetTokenCookie(c, token)
   |      (middleware/auth.go:27-36)
   |      - Sets "forge_token" cookie:
   |          HTTPOnly=true, SameSite=Lax, Secure=false, Path=/
   |
   +-- 5. Return 200 with user JSON: { id, username, email, avatar,
                                       bio, role, stacks, created_at }
         (password is never returned; see models.go:21 `json:"-"`)
```

### 3.3 Token validation on protected routes

Any route registered with `sessionManager.RequireAuth` (e.g.
`POST /api/v1/threads`, `GET /api/v1/auth/me`) is wrapped by this
middleware (middleware/auth.go:50-68):

```
RequireAuth(c)
   |
   +-- 1. extractToken(c)  (middleware/auth.go:103-114)
   |      - Prefer cookie "forge_token"
   |      - Else Authorization header "Bearer <token>"
   |      - Else -> 401 {"error":"Unauthorized"}
   |
   +-- 2. lib.VerifyJWT(token, secret)  (lib/jwt.go:30-46)
   |      - Checks signing method is HMAC (rejects alg confusion attacks)
   |      - Validates exp/iat
   |      - On any error -> 401 {"error":"Unauthorized"}
   |
   +-- 3. c.Locals("user_id", claims.UserID)
        c.Locals("user_role", claims.Role)
        c.Next()
```

Downstream handlers retrieve the identity via
`h.SessionManager.GetUserID(c)` and/or `GetUserRole(c)`
(middleware/auth.go:70-101). Both have a fallback that re-parses the
cookie if the locals were not set (e.g. when `RequireAuth` was not in the
chain).

### 3.4 Logout

```
HTTP POST /api/v1/auth/signout  (RequireAuth)
   |
   v
authHandler.SignoutHandler  (handlers/auth.go:89-94)
   |
   +-- 1. h.SessionManager.ClearTokenCookie(c)
   |      (middleware/auth.go:38-48)
   |      - Sets "forge_token" with MaxAge=-1 to instruct browser to drop it
   |
   +-- 2. Return 200 {"message":"Logged out successfully"}
```

The JWT itself is not invalidated server-side (no blacklist table); clients
simply drop the cookie. Note that any client holding the raw token could
still present it via `Authorization: Bearer` until expiry. A token blocklist
would be needed to fully revoke.

---

## 4. Request lifecycle

A typical authenticated request flows like this. `GET /api/v1/threads/:slug`
is **not** authenticated, so we use `POST /api/v1/threads` instead to show
the full chain:

```
Client
  |  HTTP POST /api/v1/threads
  |  Cookie: forge_token=eyJhbGciOiJIUzI1NiIs...
  |  Body: {"title":"...","content":"...","tags":["go"],"status":"draft"}
  v
Fiber router
  |  1. recover  (catches panic from any later step)
  |  2. logger   (logs method, path, status, latency)
  |  3. cors     (adds Access-Control-* headers, handles OPTIONS preflight)
  |  4. (no authLimiter on this route)
  |  5. sessionManager.RequireAuth
  |       - extracts token from cookie
  |       - verifies signature + expiry
  |       - sets c.Locals("user_id", ...), c.Locals("user_role", ...)
  |  6. threadHandler.CreateThreadHandler
  v
Handler (handlers/thread.go:35-58)
  |  - c.BodyParser(&req)
  |  - userID := h.SessionManager.GetUserID(c)
  |  - h.ThreadService.Create(ctx, title, content, status, tags, userID)
  v
Use case (usecase/thread_service.go:25-60)
  |  - validates lengths (title 5-255, content 10-50000, status enum)
  |  - repo.GenerateUniqueSlug(ctx, title)  (lib/slug.go:30-67)
  |  - builds *domain.Thread
  |  - repo.Create(ctx, thread, tags)        (adapters/db/thread_repository.go:22-48)
  |      - find existing tag rows by name
  |      - INSERT thread, then INSERT m2m rows in thread_tags
  |      - Preload("Author", "Tags") and map back to domain
  v
Database (PostgreSQL 16)
  |  - Returns inserted row + joined author + tags
  v
Use case
  |  - returns *domain.Thread
  v
Handler
  |  - mapThreadToResponse(thread)  (handlers/thread.go:210-229)
  |  - returns c.Status(201).JSON(...)
  v
Fiber
  |  - logger records the response status
  v
Client
  |  HTTP 201 Created
  |  {"id":1, "title":"...", "slug":"...", "author":{...}, "tags":[{...}], ...}
```

### Error short-circuits

Any step can short-circuit by returning a non-nil `error` from the handler.
Fiber's `ErrorHandler` (set in `main.go:34-37`) will then return
`500 {"error": "..."}`. Handlers in this codebase generally pre-check
errors themselves and return appropriate 4xx codes via
`c.Status(...).JSON(fiber.Map{"error": ...})`, so the global ErrorHandler
mostly fires for unexpected runtime errors (e.g. DB outages).

---


## 5. CRUD workflows

### 5.1 Threads

#### Create  (`POST /api/v1/threads`, RequireAuth)

```
handler  (handlers/thread.go:35-58)
  +-- BodyParser -> CreateThreadRequest{title, content, tags[], status}
  +-- userID := SessionManager.GetUserID(c)
  +-- ThreadService.Create(ctx, title, content, status, tags, userID)
  |    (usecase/thread_service.go:25-60)
  |    +-- validate title length 5-255  -> ErrInvalidInput (handler->400)
  |    +-- validate content length 10-50000
  |    +-- default status="draft" if empty
  |    +-- validate status in {draft, published}
  |    +-- repo.GenerateUniqueSlug(ctx, title)
  |    |   (lib/slug.go:30-67)
  |    |   - GenerateSlug: lowercase, [^a-z0-9]+ -> "-", collapse "--", trim "-"
  |    |   - Fall back to "thread" if empty
  |    |   - Try base, base-1, ..., base-10
  |    |   - On miss: append rand(1000-9999) or rand(10000-99999)
  |    +-- build *domain.Thread
  |    +-- repo.Create(ctx, thread, tagNames)
  |        (adapters/db/thread_repository.go:22-48)
  |        - select existing Tag rows whose name is in tagNames
  |        - INSERT thread row
  |        - Preload Author + Tags, map to domain
  +-- 201 + mapThreadToResponse
```

#### List  (`GET /api/v1/threads`, public)

Query params: `page` (default 1), `pageSize` (default 5, clamped to 50),
`sort` (`latest` (default) | `oldest` | `votes`).

```
handler  (handlers/thread.go:60-88)
  +-- ThreadService.List(ctx, page, pageSize, sort)
       (usecase/thread_service.go:62-73) -- clamps page/pageSize
       +-- repo.List(ctx, page, pageSize, sort)
            (adapters/db/thread_repository.go:50-85)
            - COUNT(*) WHERE status='published' (total)
            - SELECT threads.*, subqueries for upvotes/downvotes/replies_count
            - Order by created_at DESC | ASC | by net votes DESC
            - Preload Author, Tags
            - Offset + Limit
            - Map each row via threadFromModel
       Return: ([]domain.Thread, total, nil)
  +-- Build pagination {page, pageSize, total, totalPages=ceil(total/pageSize)}
  +-- 200 { threads: [...], pagination: {...} }
```

#### Featured  (`GET /api/v1/threads/featured`, public)

```
handler  (handlers/thread.go:90-99)
  +-- ThreadService.GetFeatured(ctx)
       +-- repo.GetFeatured(ctx)
            (adapters/db/thread_repository.go:120-140)
            - Window: threads created within the last 7 days
            - Status = 'published'
            - Order by SUM(votes.value) DESC
            - First row only
            - 404 if no row found (handler->404)
```

#### Trending  (`GET /api/v1/threads/trending`, public)

```
Same as Featured but LIMIT 3 and no time window.
(handlers/thread.go:101-117, adapters/db/thread_repository.go:142-167)
```

#### Read by slug  (`GET /api/v1/threads/:slug`, public)

```
handler  (handlers/thread.go:119-139)
  +-- slug := c.Params("slug")
  +-- ThreadService.GetBySlug(ctx, slug)
  |    (usecase/thread_service.go:96-106)
  |    +-- repo.GetBySlug(ctx, slug) -- Preload Author, Tags,
  |    |     Comments where parent_id IS NULL, Comments.Replies,
  |    |     Comments.Author, Comments.Replies.Author
  |    |     (adapters/db/thread_repository.go:169-184)
  |    +-- repo.IncrementViewCount(ctx, thread.ID)  -- best-effort
  |    +-- thread.ViewCount++  (in-memory)
  +-- If thread.Status != "published" -> 404 (hides drafts)
  +-- 200 + thread + nested comments (via serializeComments)
```

#### Update  (`PATCH /api/v1/threads/:slug`, RequireAuth, owner-only)

```
handler  (handlers/thread.go:148-182)
  +-- BodyParser -> UpdateThreadRequest{title?, content?, status?, tags?}
  |   (pointer fields: nil = "do not change")
  +-- ThreadService.Update(ctx, slug, userID, title*, content*, status*, tags)
  |    (usecase/thread_service.go:108-154)
  |    +-- GetBySlug -> ErrThreadNotFound (handler->404)
  |    +-- if thread.AuthorID != userID -> ErrPermissionDenied (handler->403)
  |    +-- If title is being changed:
  |    |     - validate length
  |    |     - if title differs (case-insensitive), regenerate slug
  |    +-- If content is being changed: validate length 10-50000
  |    +-- If status is being changed: must be draft|published
  |    +-- repo.Update(ctx, thread, tagNames)
  |         (adapters/db/thread_repository.go:186-217)
  |         - Save() thread row
  |         - if tagNames != nil: Association("Tags").Replace(&tags)
  |           (GORM deletes old m2m rows, inserts new ones)
  |         - Re-preload and map to domain
  +-- 200 + updated thread
```

#### Delete  (`DELETE /api/v1/threads/:slug`, RequireAuth, owner-only)

```
handler  (handlers/thread.go:184-208)
  +-- ThreadService.Delete(ctx, slug, userID)
  |    +-- GetBySlug -> ErrThreadNotFound (->404)
  |    +-- if thread.AuthorID != userID -> ErrPermissionDenied (->403)
  |    +-- repo.Delete(ctx, thread)
  |         (adapters/db/thread_repository.go:219-221)
  |         - r.db.Delete(&models.Thread{}, t.ID)  -- soft delete (DeletedAt set)
  +-- 200 {"message":"Thread deleted"}
```

### 5.2 Comments

#### Create  (`POST /api/v1/threads/:slug/comments`, RequireAuth)

Body: `{ content, parentId? }` (`parentId` is `*uint`).

```
handler  (handlers/comment.go:31-68)
  +-- BodyParser
  +-- slug := c.Params("slug")
  +-- userID := SessionManager.GetUserID(c)
  +-- CommentService.Create(ctx, slug, content, parentID, userID)
  |    (usecase/comment_service.go:28-65)
  |    +-- validate content length 2-10000 -> ErrInvalidInput (->400)
  |    +-- threadRepo.GetBySlug(ctx, slug) -> ErrThreadNotFound (->404) if missing
  |    +-- if parentID != nil:
  |    |     - repo.GetByID(parentID) -> ErrCommentNotFound (->400) if missing
  |    |     - parent.ThreadID must match thread.ID -> ErrInvalidInput (->400)
  |    |     - parent.ParentID must be nil -> ErrInvalidInput (->400)
  |    |       (enforces 1-level nesting only)
  |    +-- build *domain.Comment
  |    +-- repo.Create(ctx, comment)
  |         (adapters/db/comment_repository.go:20-39)
  |         - INSERT, Preload Author, map to domain
  +-- 201 { message, comment }
```

#### Read (nested)  (`GET /api/v1/threads/:slug`)

Comments are returned as part of the thread detail. The repository uses
GORM preloads to fetch top-level comments and their replies in two
queries:

```
Preload("Comments", "parent_id IS NULL")   // top-level only
Preload("Comments.Replies")                // for each top-level
Preload("Comments.Author")
Preload("Comments.Replies.Author")
```

`serializeComments` (handlers/thread.go:243-275) then nests the reply
slice under each parent.

#### Delete  (`DELETE /api/v1/comments/:id`, RequireAuth, owner or admin)

```
handler  (handlers/comment.go:70-101)
  +-- id := c.ParamsInt("id")
  +-- userID := SessionManager.GetUserID(c)
  +-- userRole := SessionManager.GetUserRole(c)
  +-- CommentService.Delete(ctx, id, userID, userRole)
  |    (usecase/comment_service.go:67-83)
  |    +-- repo.GetByID -> ErrCommentNotFound (->404)
  |    +-- if comment.AuthorID != userID && userRole != "admin"
  |    |     -> ErrCommentForbidden (->403)
  |    +-- repo.DeleteReplies(ctx, comment.ID) -- cascade-delete children
  |    +-- repo.Delete(ctx, comment.ID) -- soft delete parent
  +-- 200 {"message":"Comment deleted"}
```

### 5.3 Votes

#### Upvote / Downvote / Retract  (RequireAuth)

```
POST /api/v1/threads/:slug/vote
POST /api/v1/comments/:id/vote
Body: { "value": -1 | 0 | 1 }
```

```
handler  (handlers/vote.go:30-103)
  +-- BodyParser -> VoteRequest{value}
  +-- For thread route: slug := c.Params("slug"); userID := SessionManager.GetUserID(c)
  |  For comment route: id := c.ParamsInt("id"); userID := SessionManager.GetUserID(c)
  +-- VoteService.VoteThread(ctx, slug, userID, value) or VoteComment(...)
  |    (usecase/vote_service.go:25-57)
  |    +-- validate value in {-1, 0, 1} -> ErrInvalidInput (->400)
  |    +-- resolve target: threadRepo.GetBySlug or commentRepo.GetByID
  |    |     - thread not found -> ErrThreadNotFound (->404)
  |    |     - comment not found -> ErrCommentNotFound (->404)
  |    +-- repo.VoteThread/VoteComment(ctx, targetID, userID, value)
  |         (adapters/db/vote_repository.go:20-48)
  |         - If value == 0: DELETE WHERE user_id=? AND thread_id=?
  |           (retract the vote)
  |         - Else: UPSERT via clause.OnConflict{user_id, thread_id}
  |           -> DoUpdates SET value=EXCLUDED.value
  |         - PostgreSQL ON CONFLICT means one row per (user, target)
  +-- repo.GetThreadVotes/GetCommentVotes -> (upvotes, downvotes)
  +-- 200 { message, upvotes, downvotes }
```

#### Counting votes

`thread_repository.go:50-85` uses correlated subqueries in `SELECT` to
compute `upvotes`, `downvotes`, and `replies_count` per thread in a single
query:

```sql
SELECT threads.*,
  (SELECT COUNT(*) FROM votes WHERE votes.thread_id = threads.id
     AND votes.value = 1)   AS upvotes,
  (SELECT COUNT(*) FROM votes WHERE votes.thread_id = threads.id
     AND votes.value = -1)  AS downvotes,
  (SELECT COUNT(*) FROM comments WHERE comments.thread_id = threads.id
     AND comments.parent_id IS NULL) AS replies_count
FROM threads WHERE threads.status = 'published'
```

For comments, `commentFromModel` calls `m.Upvotes(db)` and `m.Downvotes(db)`
(models.go:128-138) for each comment when serializing a thread.

### 5.4 User profile

#### Get  (`GET /api/v1/users/:username`, public)

```
handler  (handlers/user.go:52-66)
  +-- username := c.Params("username")
  +-- UserService.GetUserProfile(ctx, username)
  |    (usecase/user_service.go:23-26) -- thin pass-through
  |    +-- repo.GetByUsername(ctx, username)  (adapters/db/user_repository.go:56-65)
  |         - errors.Is(err, gorm.ErrRecordNotFound) -> domain.ErrNotFound
  |         - fromModel maps models.User -> domain.User (parses Stacks JSON)
  +-- 200 { user: { id, username, avatar, bio, stacks, role, created_at } }
       -> 404 if not found
```

#### Update  (`PATCH /api/v1/users/:username`, RequireAuth, self only)

```
handler  (handlers/user.go:69-134)
  +-- username := c.Params("username")
  +-- userID := SessionManager.GetUserID(c)  -> 401 if 0
  +-- BodyParser -> UpdateUserRequest{avatar?, bio?, stacks?}
  |   (pointer fields distinguish "absent" from "empty")
  +-- Validate lengths:
  |     - bio <= 500 chars -> 400
  |     - stacks <= 10 items -> 400
  +-- Fetch user via GetUserProfile -> 404 if missing
  +-- Security: if user.ID != userID -> 403 (forbid editing others)
  +-- UserService.UpdateProfile(ctx, userID, updates)
  |    (usecase/user_service.go:29-54)
  |    +-- repo.GetByID(ctx, userID)
  |    +-- If updates.Bio != "" -> user.Bio = updates.Bio
  |    +-- If updates.Avatar != "" -> user.Avatar = updates.Avatar
  |    +-- If updates.Stacks != nil -> user.Stacks = updates.Stacks
  |    +-- repo.Update(ctx, user)
  |         (adapters/db/user_repository.go:80-103)
  |         - Update with Omit("CreatedAt"), Select("*")
  |         - returns domain.ErrNotFound if RowsAffected == 0
  +-- 200 { user: <updated> }
```

#### User's threads  (`GET /api/v1/users/:username/threads`, public)

```
handler  (handlers/user.go:137-203)
  +-- GetUserProfile -> 404 if user missing
  +-- ThreadService.ListByUser(ctx, username, page, pageSize)
  |    (usecase/thread_service.go:75-86 -> repo.ListByUser)
  |    - Find user by username -> 404 (or empty result)
  |    - COUNT threads WHERE author_id = user.ID
  |    - SELECT threads with vote/comment subqueries, ordered by created_at DESC
  |    - Preload Author, Tags
  +-- 200 { threads, pagination: { page, pageSize, totalItems, totalPages } }
```

### 5.5 Tags

#### List  (`GET /api/v1/tags`, public)

```
handler  (handlers/tag.go:30-51)
  +-- TagService.ListTags(ctx)  (usecase/tag_service.go:23-25, thin pass-through)
  +-- repo.ListTags(ctx)
       (adapters/db/tag_repository.go:21-36)
       - SELECT * FROM tags ORDER BY name ASC
  +-- 200 { tags: [{id, name, color}, ...] }
```

#### Create  (`POST /api/v1/tags`, RequireAuth, admin only)

```
handler  (handlers/tag.go:54-94)
  +-- BodyParser -> CreateTagRequest{name, color?}
  +-- userRole := SessionManager.GetUserRole(c)
  +-- TagService.CreateTag(ctx, name, color, userRole)
  |    (usecase/tag_service.go:27-60)
  |    +-- if userRole != "admin" -> "Admin access required" (handler->403)
  |    +-- trim whitespace
  |    +-- name length 3-50  (handler->400)
  |    +-- color matches ^#[0-9a-fA-F]{6}$ if non-empty
  |    |   default to "#6366f1" if empty
  |    +-- repo.GetByName(ctx, name) -- if found -> "Tag with this name already exists"
  |    |   (handler->409)
  |    +-- repo.CreateTag(ctx, &domain.Tag{Name, Color})
  |         (adapters/db/tag_repository.go:38-50)
  +-- 201 { tag: {id, name, color} }
```

Note: the service returns plain `errors.New(...)` strings rather than typed
errors, so the handler matches on `err.Error()` text. This is brittle
(refactoring a message would break the handler) but is the current
implementation.

---

## 6. Error handling flow

Three layers participate:

### Layer 1: Domain
`internal/domain/errors.go` defines sentinel errors that the use-case layer
wraps with `fmt.Errorf("...: %w", domain.ErrXxx)`. Examples:
- `domain.ErrNotFound`
- `domain.ErrForbidden`
- `domain.ErrInvalidInput`

### Layer 2: Use cases
Each service method either returns a domain error (wrapped) or a
service-specific typed error:
- `usecase.ErrThreadNotFound`  = `fmt.Errorf("thread: %w", domain.ErrNotFound)`
- `usecase.ErrPermissionDenied` = `fmt.Errorf("thread: %w", domain.ErrForbidden)`
- `usecase.ErrCommentNotFound`, `usecase.ErrCommentForbidden`
- `usecase.ErrInvalidInput`, `ErrUsernameTaken`, `ErrEmailRegistered`, `ErrInvalidCredentials`

Handlers match these via `errors.Is` (preferred) or by string comparison
in `tag.go` (legacy pattern).

### Layer 3: Handlers
A typical handler error block:

```go
if err != nil {
    if errors.Is(err, usecase.ErrThreadNotFound) { return c.Status(404)... }
    if errors.Is(err, usecase.ErrPermissionDenied) { return c.Status(403)... }
    if errors.Is(err, domain.ErrInvalidInput)     { return c.Status(400)... }
    return c.Status(500)... // catch-all
}
```

### Global fallback
`fiber.Config{ ErrorHandler: ... }` (main.go:34-37) catches any error
that escapes a handler and returns 500 with the error message. In
practice, all handlers cover their error paths explicitly, so this mostly
fires for nil-pointer or DB-driver failures.

### Status code map

| Layer signal | HTTP status |
|---|---|
| `domain.ErrInvalidInput` / `usecase.ErrInvalidInput` / body parse fail | **400** |
| Auth missing or invalid (RequireAuth) | **401** |
| Login failed | **401** (with generic "Invalid credentials" to avoid user-enumeration) |
| `domain.ErrForbidden` / `usecase.ErrPermissionDenied` | **403** |
| Resource not found (`ErrNotFound` / `ErrThreadNotFound` / `ErrCommentNotFound`) | **404** |
| Duplicate tag | **409** |
| Anything else | **500** |

---

## 7. Cross-cutting concerns

### 7.1 Soft deletes
All domain models (except `Session`) embed `gorm.DeletedAt` with
`gorm:"index"`. GORM's default `Delete()` sets the timestamp and excludes
the row from normal queries. Hard deletes require `Unscoped().Delete()`.

### 7.2 Slug generation
`lib.GenerateSlug` lowercases the title, replaces any non-`[a-z0-9]` run
with a single hyphen, trims leading/trailing hyphens, and falls back to
`"thread"` for fully-symbolic input. `GenerateUniqueSlug` then probes
the `threads.slug` column for collisions and appends `-N` (up to `-10`)
or a random 4-5 digit suffix.

### 7.3 Vote upsert semantics
The `Vote` table's unique indexes are what make the `ON CONFLICT ... DO
UPDATE` pattern work. A user can switch their vote from +1 to -1 and back
without creating duplicate rows, and `value: 0` deletes the row entirely
to retract a vote.

### 7.4 Computed fields on Thread
`Upvotes`, `Downvotes`, and `RepliesCount` are tagged `gorm:"-"` on
`models.Thread` (models.go:68-70), meaning they are **not** columns.
They are populated at query time by:
- Correlated subqueries in the thread list / featured / trending queries
- The `commentFromModel` helper (which calls `m.Upvotes(db)` /
  `m.Downvotes(db)`) for comments embedded inside a thread response

### 7.5 JSON serialization
- Domain `User.Stacks` is `[]string`; `models.User.Stacks` is a JSONB
  string. The repository marshals/unmarshals in `toModel`/`fromModel`
  (adapters/db/user_repository.go:107-148).
- `models.User.Password` has `json:"-"` and is never included in any
  response.
- Response shapes are constructed with `fiber.Map` literals in handlers;
  there is no shared response DTO module.

### 7.6 CORS
- Allowed origins: from `cfg.CORSOrigins` (default:
  `http://localhost:3000,http://localhost:8080,http://127.0.0.1:8080`).
- Credentials allowed (so the JWT cookie can be sent cross-origin).
- Allowed methods: `GET,POST,PUT,DELETE,PATCH,OPTIONS`.
- Allowed headers: `Origin,Content-Type,Accept,Authorization`.

### 7.7 Rate limiting
Applied only to `POST /api/v1/auth/signup` and `POST /api/v1/auth/signin`
via Fiber's `limiter` middleware. Key is the client IP; limit is 10 req
per minute; `SkipSuccessfulRequests: true` means successful signups /
signins do **not** count toward the limit, so real users won't be
throttled while attackers who fail repeatedly are.

### 7.8 Logging
Fiber's default `logger` middleware writes request lines to stdout. The
application does not use structured logging anywhere else; the seeder
uses `log` for its one-time messages.

### 7.9 Context propagation
`context.Context` is passed as the first argument to every repository
method and is forwarded to GORM via `db.WithContext(ctx)`. This is the
standard way to allow request cancellation and timeouts to propagate
through the layers.

### 7.10 Hexagonal layering rules (enforced by convention)
- `domain/` imports nothing from the project except standard library.
- `ports/` imports only `domain/` and standard library.
- `usecase/` imports `domain/` and `ports/`.
- `adapters/db/` imports `domain/`, `models/`, and the chosen DB driver.
- `handlers/` imports `ports/`, `domain/` (only for helper types),
  `usecase/` (for error sentinels), and `middleware/`.
- `cmd/server/main.go` is the **only** package that knows about every
  layer.

### 7.11 Tests
- Test files live in `tests/` and use the `package` names of the files
  they test. They use `DATA-DOG/go-sqlmock` for unit-level DB mocking
  and `stretchr/testify` for assertions.
- Tests do not require a running database.

### 7.12 Known limitations
- `JWT_EXPIRY` env var is read but not honored (hard-coded 72h in
  config.go:45).
- No token revocation / blocklist -- sign-out only clears the cookie.
- No refresh tokens; users must re-authenticate after 72h.
- `Session` model is defined but unused at runtime.
- AutoMigrate is non-destructive; column drops must be done manually.
- Tag service uses string-matching on errors rather than typed errors.
- No request-ID / structured logging middleware.
- No graceful shutdown (`app.Listen` blocks until error; no signal
  handling).
