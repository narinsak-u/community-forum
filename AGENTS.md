# Agent Guidelines for Community Forum

## Project Overview

This is a full-stack community forum with React frontend and Go backend.

```
community-forum/
├── frontend/          # React + TypeScript frontend
└── backend/           # Go Fiber backend
```

---

## Frontend Commands

### Development
```bash
cd frontend
bun run dev          # Start dev server on port 8080
bun run preview      # Preview production build
```

### Build
```bash
bun run build        # Production build to dist/
bun run build:dev   # Development build
```

### Linting
```bash
bun run lint         # Run ESLint on all files
```

### Testing (Vitest)
```bash
bun test             # Run all tests once
bun run test:watch   # Run tests in watch mode

# Run a single test file
bunx vitest run frontend/src/test/example.test.ts

# Run tests matching a pattern
bunx vitest run --grep "example"

# Run tests in a specific directory
bunx vitest run frontend/src/test/
```

---

## Backend Commands

### Running the Server
```bash
cd backend

# Start PostgreSQL (required)
docker-compose up -d

# Run server
go run ./cmd/server

# Build
go build -o server ./cmd/server
```

### Testing
```bash
cd backend
go test ./...           # Run all tests
go test ./... -v        # Run with verbose output
go test -run TestName   # Run specific test
```

### Linting (Go)
```bash
cd backend
go fmt ./...           # Format code
go vet ./...           # Run go vet
```

---

## Frontend Code Style

### General
- React 18 + TypeScript + Vite with Tailwind CSS and shadcn/ui
- Use `@/` path alias for imports (configured in tsconfig.json)
- TypeScript strict mode is disabled - implicit `any` allowed

### Imports
- Use absolute imports: `import Button from "@/components/ui/button"`
- Group: external libs → internal components → local utilities
- Use `type` keyword for type-only imports

### Components
- Functional components with arrow functions or `function` declarations
- Use React.forwardRef for ref forwarding
- Use CVA (class-variance-authority) for variants
- Named exports with default for page components

### Naming
- **Components**: PascalCase (Button, AppLayout)
- **Hooks**: camelCase with `use` prefix (useToast)
- **Utilities**: camelCase (cn, formatDate)
- **Files**: PascalCase for components, camelCase for utils

### Types
- interfaces for object shapes
- type for unions/tuples
- Use `VariantProps<T>` from CVA for variant types

### CSS
- Use Tailwind utility classes
- Use `cn()` from `@/lib/utils` to merge classes
- Custom classes: `panel`, `terminal-label`, `heading-display`

### React Patterns
- TanStack Query for data fetching
- React Router v6 for routing
- react-hook-form + zod for forms

### What NOT to Do (Frontend)
- No default exports for components
- No CSS files - use Tailwind
- No creating folders at root level

---

## Backend Code Style (Go)

### General
- Go 1.24 with Fiber web framework
- GORM for ORM with PostgreSQL
- Follow standard Go project layout

### Structure
```
backend/
├── cmd/server/       # Entry point (main.go)
├── internal/
│   ├── handlers/    # HTTP handlers
│   ├── models/      # GORM models
│   ├── schemas/    # Request/response types
│   └── middleware/ # Custom middleware
├── migrations/     # DB migrations
└── docker-compose.yml
```

### Naming Conventions
- **Files**: lowercase with underscores (user_model.go)
- **Functions**: PascalCase exported, camelCase unexported
- **Constants**: UPPER_SNAKE_CASE
- **Packages**: lowercase, single word preferred

### Error Handling
- Return errors with context using `fmt.Errorf("context: %w", err)`
- Use custom error types for domain errors
- Handle errors at the handler level, return proper HTTP codes

### Database
- Use GORM for all database operations
- Follow model structure in `internal/models/models.go`
- Use migrations for schema changes

### What NOT to Do (Backend)
- No business logic in handlers - delegate to services
- Don't use global variables for state
- Don't hardcode connection strings - use env vars

---

## Project Structure

### Frontend (`frontend/src/`)
```
src/
├── components/ui/    # shadcn/ui components
├── components/forge/ # Custom app components
├── hooks/           # Custom React hooks
├── lib/             # Utilities
├── pages/           # Route pages
├── test/            # Test files
└── App.tsx          # Root component
```

### Backend (`backend/`)
```
backend/
├── cmd/server/      # Entry point
├── internal/         # Private packages
│   ├── handlers/    # HTTP handlers
│   ├── models/       # GORM models
│   └── middleware/   # Middleware
└── docker-compose.yml
```

---

## Testing

### Frontend
- Vitest + @testing-library/react
- Files: `*.test.ts` or `*.test.tsx` in `src/test/`

### Backend
- Standard Go testing
- Files: `*_test.go` alongside source files