# Backend Overview

> Documentation for the Community Forum backend service. Last updated: 2026-06-08.
>
> **Important:** This project has been refactored from the original flat
> `models/handlers/schemas` structure into a **Hexagonal Architecture**
> (Ports and Adapters). The root `README.md` still shows the old structure;
> this document reflects the current state. See
> [`docs/BACKEND-IMPROVEMENT.md`](./BACKEND-IMPROVEMENT.md) for the rationale.

---

## 1. Project Summary

A RESTful JSON API for a community forum. Users sign up, create threaded
discussions, post comments (with one level of nested replies), tag threads,
and upvote / downvote threads and comments. The service exposes a versioned
JSON API at `/api/v1/...`, issues JWT cookies for sessions, persists data in
PostgreSQL via GORM, and auto-seeds a demo dataset on first start.

**Tech stack (one-liner):** Go 1.24 + Fiber v2 + GORM + PostgreSQL 16, wired
together with explicit dependency injection in `cmd/server/main.go`.

---

## 2. Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                        # Entry point / composition root
├── internal/
│   ├── adapters/
│   │   └── db/                            # Outbound adapters (GORM)
│   │       ├── user_repository.go
│   │       ├── thread_repository.go
│   │       ├── comment_repository.go
│   │       ├── vote_repository.go
│   │       └── tag_repository.go
│   ├── config/
│   │   └── config.go                      # Env loading, DB init, auto-migrate, seed
│   ├── domain/                            # Pure domain entities (no infra deps)
│   │   ├── user.go
│   │   ├── thread.go
│   │   ├── comment.go
│   │   ├── tag.go
│   │   ├── vote.go
│   │   └── errors.go
│   ├── handlers/                          # Inbound adapters (Fiber HTTP)
│   │   ├── auth.go
│   │   ├── thread.go
│   │   ├── comment.go
│   │   ├── vote.go
│   │   ├── user.go
│   │   └── tag.go
│   ├── lib/                               # Cross-cutting helpers
│   │   ├── jwt.go
│   │   └── slug.go
│   ├── middleware/
│   │   └── auth.go                        # SessionManager (JWT cookie + RequireAuth)
│   ├── models/                            # GORM-tagged DB models (separate from domain)
│   │   └── models.go
│   ├── ports/                             # Interfaces (contracts)
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── thread.go
│   │   ├── comment.go
│   │   ├── vote.go
│   │   └── tag.go
│   ├── seed/
│   │   └── seed.go                        # First-run demo data
│   └── usecase/                           # Business logic (implements inbound ports)
│       ├── auth_service.go
│       ├── user_service.go
│       ├── thread_service.go
│       ├── comment_service.go
│       ├── vote_service.go
│       └── tag_service.go
├── tests/                                 # Black-box tests (mirror internal/ tree)
│   ├── adapters/
│   ├── domain/
│   ├── handlers/
│   ├── lib/
│   └── usecase/
├── docs/                                  # You are here
│   ├── OVERVIEW.md
│   ├── WORKFLOW.md
│   ├── BACKEND-IMPROVEMENT.md
│   ├── BACKEND-PLAN.md
│   └── plans/
├── docker-compose.yml                     # PostgreSQL + API service
├── Dockerfile                             # Multi-stage Go build -> alpine
├── .env / .env.example                    # Environment variables
├── go.mod / go.sum
└── .dockerignore
```

> **Note:** the original `migrations/` directory referenced in the root
> `README.md` does not exist. Schema is created via GORM's `AutoMigrate` in
> `internal/config/config.go:64-74`.

---

## 3. Tech Stack Details

### Go toolchain
- **Go version:** `1.24.2` (see `go.mod:3`)
- **Module path:** `community-forum/backend`

### Direct dependencies (from `go.mod`)

| Module | Version | Purpose |
|---|---|---|
| `github.com/gofiber/fiber/v2` | v2.52.0 | HTTP framework |
| `gorm.io/gorm` | v1.25.0 | ORM |
| `gorm.io/driver/postgres` | v1.5.0 | PostgreSQL driver for GORM |
| `github.com/golang-jwt/jwt/v5` | v5.3.1 (indirect) | JWT signing/verification |
| `golang.org/x/crypto` | v0.14.0 | bcrypt password hashing |
| `github.com/joho/godotenv` | v1.5.1 | `.env` file loader |
| `github.com/DATA-DOG/go-sqlmock` | v1.5.0 | DB mocking for tests |
| `github.com/stretchr/testify` | v1.8.2 | Test assertions |

### Notable indirect dependencies
- `github.com/jackc/pgx/v5` — Postgres driver used by `gorm.io/driver/postgres`.
- `github.com/valyala/fasthttp` — Underlying HTTP engine for Fiber.

---

## 4. Configuration

### Environment variables

Loaded by `internal/config/config.go:30-47` via `godotenv`. The function
silently logs and falls back to defaults if `.env` is missing.

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5433` | PostgreSQL port (note: `5433`, not `5432`, to avoid host conflicts) |
| `DB_USER` | `postgres` | DB user |
| `DB_PASSWORD` | `postgres` | DB password |
| `DB_NAME` | `community_forum` | DB name |
| `DB_SSLMODE` | `disable` | GORM SSL mode |
| `PORT` | `8080` | HTTP listen port |
| `CORS_ORIGINS` | `http://localhost:3000,http://localhost:8080,http://127.0.0.1:8080` | Allowed CORS origins (used in `main.go:52`) |
| `JWT_SECRET` | `midnight-forge-dev-secret-change-in-production` | HMAC-SHA256 signing key |
| `JWT_EXPIRY` | _(hard-coded `72h` in `config.go:45`)_ | Token lifetime (env var is read but not yet honored) |

`.env.example` mirrors the same set of variables.

### Database setup

`docker-compose.yml` defines two services:

- **`postgres`** — `postgres:16-alpine` image
  - User/password: `postgres` / `postgres`
  - DB name: `community_forum`
  - Port mapping: `5433 -> 5432`
  - Volume: `postgres_data` (persistent)
  - Healthcheck: `pg_isready -U postgres` every 10s
- **`api`** — built from the local `Dockerfile`
  - Port mapping: `8080 -> 8080`
  - Uses the docker-network hostname `postgres` for `DB_HOST`
  - `depends_on: postgres` with `condition: service_healthy`

### Dockerfile

Two-stage build:
1. `golang:1.24-alpine` builder → `go build -o /app/server ./cmd/server` (CGO disabled, linux target).
2. `alpine:3.19` runtime with `ca-certificates` + `tzdata`, exposing port 8080.

`.dockerignore` excludes `.env`, `.git`, `docs/`, `tests/`, `*.md`, and
Docker-related files from the build context.

---

## 5. Data Models Overview

Two parallel struct hierarchies exist by design:

- **`internal/models/`** — GORM-tagged structs that map 1:1 to tables.
- **`internal/domain/`** — Pure Go structs with no infrastructure
  dependencies, used by `usecase` and `ports`. Repositories translate
  between the two.

| Model (GORM) | File:line | Domain file | Brief description |
|---|---|---|---|
| `User` | `models.go:11-33` | `domain/user.go` | Forum member. Has many threads, comments, votes, sessions. |
| `Session` | `models.go:36-44` | _(no domain type)_ | Currently unused at runtime; reserved for future session table. |
| `Thread` | `models.go:47-71` | `domain/thread.go` | Discussion topic with title, slug, content, status, view count; has many comments and tags (m2m). |
| `Comment` | `models.go:74-92` | `domain/comment.go` | Reply inside a thread. Supports one level of nested replies via `ParentID` self-reference. |
| `Tag` | `models.go:95-104` | `domain/tag.go` | Category label with a hex color. Many-to-many with threads. |
| `Vote` | `models.go:107-125` | `domain/vote.go` | Upvote/downvote on either a thread or a comment (exactly one is non-null). |

### Key relationships
- `Thread.AuthorID` → `User.ID` (FK, indexed)
- `Thread.Tags` ↔ `Tag` (m2m via implicit `thread_tags` join table)
- `Comment.ThreadID` → `Thread.ID`
- `Comment.ParentID` → `Comment.ID` (self-referencing; one level of nesting)
- `Vote.UserID` + `Vote.ThreadID` composite unique index → one vote per user per thread
- `Vote.UserID` + `Vote.CommentID` composite unique index → one vote per user per comment
- All non-join tables use GORM soft deletes (`DeletedAt gorm.DeletedAt`).

### Domain errors
`internal/domain/errors.go` defines:
- `ErrNotFound`
- `ErrForbidden`
- `ErrInvalidInput`

Service-layer code wraps these with `%w` to allow `
