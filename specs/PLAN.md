# Migration Plan

Restructure the project by separating frontend and backend into two directories:

```
community-forum/
├── frontend/   # Existing React app
└── backend/    # New Go backend
```

## Tech Stacks

### Frontend
- React 18 + TypeScript
- TanStack Query
- Zustand
- Tailwind CSS
- Zod
- Bun (runtime)

### Backend
- Go
- Fiber (https://gofiber.io/)
- GORM (ORM)
- PostgreSQL (database)
- Docker + Docker Compose

## Backend Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # GORM models
│   ├── schemas/             # Request/response types
│   └── middleware/          # Middleware (auth, logging, etc.)
├── migrations/              # Database migrations
└── docker-compose.yml       # Local dev setup with PostgreSQL
```

## Shared Types

Shared TypeScript types maintained manually in `frontend/src/types/`. Backend defines the API contract.

## Migration Approach

Big-bang migration:
1. Move frontend code to `frontend/`
2. Initialize Go backend in `backend/`
3. Set up Docker Compose with PostgreSQL
4. Create API structure, models, and handlers
5. Connect frontend to backend via TanStack Query

## Features
- Implement API, Models, Schemas to match frontend design
