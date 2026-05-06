# Backend Improvement Plan: Hexagonal Architecture Refactoring

This document outlines the strategy for refactoring the current Go backend from a coupled handler-model-database structure to a **Hexagonal Architecture** (also known as Ports and Adapters).

## 1. Objectives
- **Decouple Business Logic**: Remove dependencies on Fiber and GORM from core business rules.
- **Improve Testability**: Enable unit testing of use cases using mocks for repositories.
- **Maintainability**: Create clear boundaries between infrastructure (DB, HTTP) and the domain.
- **Flexibility**: Allow swapping out the database or web framework with minimal changes to core logic.

## 2. Target Architecture
The system will be organized into four distinct layers:

### A. Domain Layer (`internal/domain`)
- **Entities**: Plain Old Go Objects (POGOs) representing the core concepts (Thread, User, Comment).
- **Logic**: Self-contained business rules (e.g., "A user can only vote once").
- **Constraint**: MUST NOT import any other internal package or external infrastructure libraries (GORM, Fiber).

### B. Ports Layer (`internal/ports`)
- **Driving Ports (Input)**: Interfaces defining what the application can do (e.g., `ThreadService`).
- **Driven Ports (Output)**: Interfaces defining what the application needs from the outside world (e.g., `ThreadRepository`, `AuthService`).

### C. Use Case Layer (`internal/usecase`)
- **Implementation**: Realization of the Driving Ports.
- **Responsibility**: Orchestrates the flow of data between Domain entities and Driven Ports (Repositories).
- **Validation**: Performs application-level validation.

### D. Adapters Layer (`internal/adapters`)
- **Driving Adapters (`internal/adapters/http`)**: Fiber handlers that translate HTTP requests into Use Case calls.
- **Driven Adapters (`internal/adapters/db`)**: GORM implementations of the Repository interfaces.

## 3. Directory Structure Changes
```text
backend/
├── cmd/server/          # Entry point (DI setup)
└── internal/
    ├── domain/          # Entities (Thread, User, Comment)
    ├── ports/           # Interfaces (Service, Repository)
    ├── usecase/         # Business logic implementations
    └── adapters/
        ├── http/        # Fiber Handlers & Middlewares
        │   └── dto/     # Request/Response structures
        └── db/          # GORM Repositories
```

## 4. Key Refactoring Steps

### Step 1: Extract Domain Entities
- Move structs from `internal/models` to `internal/domain`.
- Remove GORM tags and database-coupled methods (like `Upvotes(db *gorm.DB)`).

### Step 2: Define Ports (Interfaces)
- Define `ThreadRepository` with methods like `Save`, `GetBySlug`, `GetFeatured`.
- Define `ThreadService` with methods like `CreateThread`, `ListThreads`.

### Step 3: Implement Repositories (Adapters)
- Create GORM-based implementations of the repository interfaces in `internal/adapters/db`.
- Move SQL-heavy logic from handlers to these repositories.

### Step 4: Implement Use Cases
- Create service implementations in `internal/usecase`.
- Move business logic (validation, slug generation, permission checks) from handlers to services.

### Step 5: Refactor Handlers
- Update Fiber handlers to depend on `ports.ThreadService` instead of `*gorm.DB`.
- Use DTOs for request/response mapping to keep domain entities clean.

## 5. Implementation Phasing
1. **Phase 1: Setup & Domain**: Create the new directory structure and migrate the `User` and `Auth` logic first as a template.
2. **Phase 2: Thread & Voting**: Migrate the core discussion logic.
3. **Phase 3: Cleanup**: Remove the old `internal/handlers` and `internal/models` once all routes are migrated.

## 6. Testing Strategy
- **Unit Tests**: Test `internal/usecase` using mock repositories (generated via `mockery`).
- **Integration Tests**: Test `internal/adapters/db` against a real PostgreSQL instance (using testcontainers).
- **E2E Tests**: Test the full HTTP flow using `httptest` and Fiber's `Test` method.
