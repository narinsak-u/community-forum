# Community Forum

A full-stack community forum application with React frontend and Go backend.

## Project Structure

```
community-forum/
├── frontend/          # React + TypeScript frontend
└── backend/           # Go Fiber backend
```

## Frontend

**Tech Stack:**

- React 18 + TypeScript
- Vite (build tool)
- Tailwind CSS
- shadcn/ui components
- TanStack Query (data fetching)
- React Router v6
- Zod (validation)
- Bun (runtime)

**Pages:**

- `/` - Home (thread list)
- `/login` - Login/Signup
- `/thread/:slug` - Thread detail
- `/create` - Create new thread
- `/profile` - User profile
- `/settings` - Settings

**Commands:**

```bash
cd frontend
bun install          # Install dependencies
bun run dev          # Start dev server (port 8080)
bun run build        # Production build
bun test             # Run tests
bun run lint         # Run ESLint
```

## Backend

**Tech Stack:**

- Go 1.24
- Fiber (HTTP framework)
- GORM (ORM)
- PostgreSQL (database)
- Docker + Docker Compose

**Structure:**

```
backend/
├── cmd/server/       # Entry point
├── internal/
│   ├── models/       # GORM models
│   ├── handlers/    # HTTP handlers
│   ├── schemas/     # Request/response types
│   └── middleware/  # Custom middleware
├── migrations/      # Database migrations
└── docker-compose.yml
```

**Models:**

- User
- Thread
- Comment
- Tag
- Vote

**Commands:**

```bash
cd backend

# Start PostgreSQL
docker-compose up -d

# Run server
go run ./cmd/server

# Build
go build -o server ./cmd/server
```

**Environment Variables:**
Copy `.env.example` to `.env`:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=community_forum
DB_SSLMODE=disable
PORT=8080
```

**API Endpoints:**

- `GET /` - Health check
- `GET /health` - Health check
- `GET /api/v1/threads` - List all threads
- `GET /api/v1/threads/:slug` - Get thread by slug

## Getting Started

1. **Frontend:**

   ```bash
   cd frontend
   bun install
   bun run dev
   ```

2. **Backend:**
   ```bash
   cd backend
   docker-compose up -d
   go run ./cmd/server
   ```

The frontend runs on `http://localhost:8080` and connects to the backend API.
