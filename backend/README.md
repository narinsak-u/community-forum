# Community Forum Backend

Go + Fiber REST API with PostgreSQL for a community forum application.

## Tech Stack

- **Go 1.24** — Runtime
- **Fiber v2** — HTTP framework
- **WebSocket** (`gofiber/contrib/websocket` + `gorilla/websocket`) — Real-time chat
- **GORM** — ORM with PostgreSQL driver
- **PostgreSQL 16** — Database
- **JWT** (`golang-jwt/jwt/v5`) — Authentication
- **bcrypt** — Password hashing

> See [`docs/OVERVIEW.md`](./docs/OVERVIEW.md) for the full architecture overview.

## Quick Start

```bash
cd backend

# Start PostgreSQL (Docker)
docker-compose up -d

# Copy environment config
cp .env.example .env

# Run server
go run ./cmd/server
```

The API starts at `http://localhost:8080`.

## Project Structure

```
backend/
├── cmd/server/main.go       # Entry point + dependency wiring (composition root)
├── internal/
│   ├── domain/             # Pure domain entities (no framework deps)
│   ├── models/             # GORM-tagged DB models
│   ├── ports/              # Interfaces (repository + service contracts)
│   ├── usecase/            # Business logic (implements service interfaces)
│   ├── adapters/db/        # GORM repository implementations
│   ├── handlers/           # HTTP + WebSocket handlers
│   ├── middleware/          # JWT auth (SessionManager)
│   ├── config/             # Env loading, DB init, auto-migrate, seed
│   ├── lib/                # JWT signing + slug generation
│   └── seed/               # Demo data seeder (idempotent)
├── tests/                   # Black-box tests (mirrors internal/)
├── docs/                    # Architecture + workflow documentation
├── docker-compose.yml       # PostgreSQL + API service
├── Dockerfile               # Multi-stage Go build → alpine
├── .env.example             # Environment variable template
├── go.mod / go.sum
└── .dockerignore
```

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters):

```
Presentation        Business Logic        Persistence
(handlers/)    →    (usecase/)       ←    (adapters/db/)
    ↓                    ↓                    ↓
   Fiber              ports/               GORM
   WebSocket        interfaces           PostgreSQL
```

See [`docs/WORKFLOW.md`](./docs/WORKFLOW.md) for detailed request lifecycle flows.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5433` | PostgreSQL port |
| `DB_USER` | `postgres` | DB user |
| `DB_PASSWORD` | `postgres` | DB password |
| `DB_NAME` | `community_forum` | DB name |
| `DB_SSLMODE` | `disable` | GORM SSL mode |
| `PORT` | `8080` | HTTP listen port |
| `CORS_ORIGINS` | `http://localhost:3000,...` | Allowed CORS origins |
| `JWT_SECRET` | *(dev secret)* | HMAC-SHA256 signing key |

## API Endpoints

### Authentication

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/auth/signup` | — | Register new user |
| POST | `/api/v1/auth/signin` | — | Login |
| POST | `/api/v1/auth/signout` | Required | Logout |
| GET | `/api/v1/auth/me` | Required | Get current user |

### Threads

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/threads` | — | List threads (paginated, sorted) |
| GET | `/api/v1/threads/featured` | — | Featured thread |
| GET | `/api/v1/threads/trending` | — | Trending threads |
| GET | `/api/v1/threads/:slug` | — | Get thread by slug |
| POST | `/api/v1/threads` | Required | Create thread |
| PATCH | `/api/v1/threads/:slug` | Required | Update thread (owner) |
| DELETE | `/api/v1/threads/:slug` | Required | Delete thread (owner) |

### Comments

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/threads/:slug/comments` | Required | Create comment |
| DELETE | `/api/v1/comments/:id` | Required | Delete comment (owner/admin) |

### Votes

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/threads/:slug/vote` | Required | Vote on thread (-1/0/+1) |
| POST | `/api/v1/comments/:id/vote` | Required | Vote on comment |

### Users

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/users` | — | List all users |
| GET | `/api/v1/users/:username` | — | Get user profile |
| PATCH | `/api/v1/users/:username` | Required | Update profile (self) |
| GET | `/api/v1/users/:username/threads` | — | Get user's threads |

### Tags

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/tags` | — | List all tags |
| POST | `/api/v1/tags` | Required | Create tag (admin) |

### Chat (WebSocket)

| Path | Protocol | Auth | Description |
|---|---|---|---|
| `/ws/chat` | WebSocket | Cookie | Real-time group chat |

WebSocket message types: `init`, `message`, `user_joined`, `user_left`, `load_more`.

## Data Models

| Model | Table | Description |
|---|---|---|
| User | `users` | Forum member |
| Thread | `threads` | Discussion topic |
| Comment | `comments` | Reply (1-level nesting) |
| Tag | `tags` | Category label |
| Vote | `votes` | Upvote/downvote on thread or comment |
| ChatMessage | `chat_messages` | Chat message in global room |

## Testing

```bash
go test ./...          # Run all tests
go test ./... -v       # Verbose
go test -run TestName  # Specific test
```

## Docs

- [`docs/OVERVIEW.md`](./docs/OVERVIEW.md) — Architecture, structure, models, configuration
- [`docs/WORKFLOW.md`](./docs/WORKFLOW.md) — Step-by-step request lifecycle flows
- [`docs/BACKEND-IMPROVEMENT.md`](./docs/BACKEND-IMPROVEMENT.md) — Hexagonal refactor rationale
