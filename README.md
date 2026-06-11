# The Lands Between

A full-stack community forum with real-time chat. React (Next.js) frontend + Go (Fiber) backend.

## Project Structure

```
├── frontend/          # Next.js 15 + React 19 + TypeScript
└── backend/           # Go 1.24 + Fiber + PostgreSQL
```

## Frontend

**Tech Stack:** Next.js 15, React 19, TypeScript, Tailwind CSS, shadcn/ui, TanStack Query

**Pages:**

| Route | Description |
|---|---|
| `/` | Landing page |
| `/threads` | Thread list (latest, popular, my posts) |
| `/thread/[slug]` | Thread detail + comments |
| `/thread/[slug]/edit` | Edit thread |
| `/create` | New thread |
| `/discussions` | Real-time group chat |
| `/network` | User directory |
| `/profile/[username]` | User profile |
| `/settings` | Account settings |
| `/login` | Sign in / Sign up |

**Commands:**

```bash
cd frontend
bun install          # Install dependencies
bun run dev          # Start dev server (port 3000/8080)
bun run build        # Production build
bun test             # Run tests
bun run lint         # Run ESLint
```

## Backend

**Tech Stack:** Go 1.24, Fiber v2, WebSocket (gorilla/websocket), GORM, PostgreSQL 16

**Architecture:** Hexagonal (Ports and Adapters) — domain / ports / usecase / adapters / handlers

**API Endpoints:**

| Category | Base Path | Auth |
|---|---|---|
| Auth | `/api/v1/auth/*` | Mixed |
| Threads | `/api/v1/threads/*` | Mixed |
| Comments | `/api/v1/threads/:slug/comments`, `/api/v1/comments/:id` | Required |
| Votes | `/api/v1/threads/:slug/vote`, `/api/v1/comments/:id/vote` | Required |
| Users | `/api/v1/users/*` | Mixed |
| Tags | `/api/v1/tags` | Mixed |
| Chat | `/ws/chat` (WebSocket) | Cookie |

**Commands:**

```bash
cd backend
docker-compose up -d  # Start PostgreSQL
go run ./cmd/server   # Start API (port 8080)
go test ./...         # Run tests
```

## Environment Variables

### Backend (`.env`)

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=community_forum
DB_SSLMODE=disable
PORT=8080
JWT_SECRET=your-secret-here
```

### Frontend (`NEXT_PUBLIC_API_URL`)

Set `NEXT_PUBLIC_API_URL=http://localhost:8080` (default) to point to the backend.

## Getting Started

1. **Start PostgreSQL:**
   ```bash
   cd backend && docker-compose up -d
   ```

2. **Start Backend:**
   ```bash
   cd backend && go run ./cmd/server
   ```

3. **Start Frontend:**
   ```bash
   cd frontend && bun install && bun run dev
   ```

Open `http://localhost:3000` (frontend). The API runs on `http://localhost:8080`.
