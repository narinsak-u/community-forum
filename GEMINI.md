# Community Forum - Project Context

A full-stack community forum application with a Go backend and a React frontend.

## Project Overview

- **Purpose**: A modern community forum platform for discussions, threading, and user interactions.
- **Backend Architecture**: Go 1.24 using the **Fiber** web framework and **GORM** for ORM. It follows a clean structure with handlers, models, and schemas.
- **Frontend Architecture**: **React 18** with **TypeScript**, built with **Vite**. It uses **Tailwind CSS** and **shadcn/ui** for styling, **TanStack Query** for data fetching, and **React Router v6** for navigation.
- **Database**: **PostgreSQL**, managed locally via **Docker Compose**.

## Directory Structure

- `backend/`: Go source code, database migrations, and Docker configuration.
  - `cmd/server/main.go`: Application entry point and API route definitions.
  - `internal/models/`: GORM database models (User, Thread, Comment, Tag, Vote).
  - `internal/handlers/`: HTTP request handlers.
  - `internal/schemas/`: Input validation and response structures.
- `frontend/`: React application.
  - `src/pages/`: Main application views (Index, ThreadDetail, Login, etc.).
  - `src/components/`: Reusable UI components (including shadcn/ui primitives).
  - `src/hooks/`: Custom React hooks.
- `specs/`: Project specifications and implementation plans.

## Building and Running

### Prerequisites
- [Go](https://go.dev/) 1.24+
- [Bun](https://bun.sh/) (Runtime for frontend)
- [Docker](https://www.docker.com/) & Docker Compose

### Backend
1. Navigate to the backend directory: `cd backend`
2. Start the database: `docker-compose up -d`
3. Configure environment variables: Copy `.env.example` to `.env` and adjust if necessary.
4. Run the server: `go run ./cmd/server` (Starts on `http://localhost:8080` by default).

### Frontend
1. Navigate to the frontend directory: `cd frontend`
2. Install dependencies: `bun install`
3. Start the development server: `bun run dev` (Runs on `http://localhost:8080` or next available port; proxying to backend might be required depending on Vite config).
4. Run tests: `bun test`

## Development Conventions

- **Backend**:
  - Uses Fiber's routing and middleware (CORS, Logger, Recover).
  - GORM AutoMigrate is used for schema management in development.
  - Models use JSON tags for consistent API responses.
- **Frontend**:
  - Functional components with TypeScript interfaces for props.
  - **shadcn/ui** components are located in `src/components/ui`.
  - **TanStack Query** is preferred for all server state management.
  - **Zod** is used for client-side and form validation.
- **General**:
  - Ensure `.env` files are never committed to version control.
  - Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and standard React best practices.
