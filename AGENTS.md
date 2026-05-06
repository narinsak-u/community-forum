# Agent Guidelines for Community Forum

## Project Overview

Full-stack community forum: React frontend + Go Fiber backend.

```
community-forum/
├── frontend/          # React + TypeScript + Vite
└── backend/           # Go Fiber + GORM + PostgreSQL
```

---

## Frontend Commands

### Development
```bash
cd frontend
bun run dev          # Dev server on port 8080
bun run preview      # Preview production build
```

### Build & Lint
```bash
bun run build        # Production build to dist/
bun run build:dev   # Development build
bun run lint        # Run ESLint
```

### Testing (Vitest)
```bash
bun test             # Run all tests once
bun run test:watch   # Watch mode

# Single test file
bunx vitest run frontend/src/test/example.test.ts

# Tests matching pattern
bunx vitest run --grep "example"

# Tests in directory
bunx vitest run frontend/src/test/
```

---

## Backend Commands

### Running Server
```bash
cd backend
docker-compose up -d    # Start PostgreSQL
go run ./cmd/server    # Run server
go build -o server ./cmd/server  # Build
```

### Testing & Linting
```bash
cd backend
go test ./...          # Run all tests
go test ./... -v       # Verbose output
go test -run TestName  # Run specific test
go fmt ./...           # Format code
go vet ./...           # Run go vet
```

---

## Frontend Code Style

### General
- React 18 + TypeScript + Vite + Tailwind CSS + shadcn/ui
- Use `@/` path alias for imports (configured in tsconfig.json)
- TypeScript strict mode disabled - implicit `any` allowed

### Imports & Naming
- Absolute imports: `import Button from "@/components/ui/button"`
- Group: external libs → internal components → local utilities
- Use `type` keyword for type-only imports
- **Components**: PascalCase (Button, AppLayout)
- **Hooks**: camelCase with `use` prefix (useToast)
- **Utilities**: camelCase (cn, formatDate)
- **Files**: PascalCase for components, camelCase for utils

### Components
- Functional components with arrow functions or `function` declarations
- Use React.forwardRef for ref forwarding
- Use CVA (class-variance-authority) for variants
- Named exports with default for page components
- interfaces for object shapes, type for unions/tuples

### CSS & Patterns
- Tailwind utility classes + `cn()` from `@/lib/utils` for class merging
- TanStack Query for data fetching
- React Router v6 for routing
- react-hook-form + zod for forms

### What NOT to Do
- No default exports for components
- No CSS files - use Tailwind
- No creating folders at root level

---

## Backend Code Style (Go)

### Structure
```
backend/
├── cmd/server/       # Entry point (main.go)
├── internal/
│   ├── handlers/    # HTTP handlers
│   ├── models/      # GORM models
│   ├── schemas/    # Request/response types
│   └── middleware/ # Custom middleware
├── migrations/      # DB migrations
└── docker-compose.yml
```

### Naming & Error Handling
- **Files**: lowercase with underscores (user_model.go)
- **Functions**: PascalCase exported, camelCase unexported
- **Constants**: UPPER_SNAKE_CASE
- **Packages**: lowercase, single word preferred
- Return errors with context: `fmt.Errorf("context: %w", err)`
- Use custom error types for domain errors
- Handle errors at handler level, return proper HTTP codes

### Database
- Use GORM for all database operations
- Follow model structure in `internal/models/models.go`
- Use migrations for schema changes

### What NOT to Do
- No business logic in handlers - delegate to services
- Don't use global variables for state
- Don't hardcode connection strings - use env vars

---

## Testing

- **Frontend**: Vitest + @testing-library/react, files in `src/test/*.test.ts`
- **Backend**: Standard Go testing, files `*_test.go` alongside source files