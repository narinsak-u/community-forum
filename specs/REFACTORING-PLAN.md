# Refactoring Plan: Codebase Architecture Improvement

## Overview

This document captures architectural friction points and deepening opportunities identified through codebase exploration. The goal is to transform shallow modules into deep modules using hexagonal architecture (backend) and React best practices (frontend).

---

## Candidates for Module Deepening

### 1. Backend: Complete Hexagonal Service Layer Migration

**Cluster**: `handlers/thread.go`, `handlers/comment.go`, `handlers/vote.go`, `handlers/tag.go` → service ports  
**Why they're coupled**: All directly use `*gorm.DB` (infrastructure leak). `auth.go` handler already uses ports correctly—these 4 should follow the same pattern.  
**Dependency category**: **Infrastructure leakage** — DB access in HTTP layer violates hexagonal architecture  
**Test impact**: Replace current handler tests that require DB with mock-based unit tests

**Architecture Pattern**: Hexagonal (Ports & Adapters)

- **Port**: Define service interface (e.g., `ThreadService` in `ports/thread.go`)
- **Adapter**: Implement in `adapters/db/thread_repository.go`
- **Current**: Handlers inject `*gorm.DB` directly
- **Target**: Handlers inject service port interface

---

### 2. Backend: Injectable Session Store (Replace Global)

**Cluster**: `middleware/auth.go` global `var Store *session.Store`  
**Why it's coupled**: Hidden global state makes testing impossible without DB. All handlers share this implicitly via package.  
**Dependency category**: **Global state** — tight coupling via package-level variable  
**Test impact**: Current tests can't test auth middleware in isolation

**Architecture Pattern**: Dependency Injection

- **Current**: Global variable `var Store *session.Store`
- **Target**: Constructor injection via handler struct field

---

### 3. Frontend: Unified Thread Data Orchestration

**Cluster**: `hooks/use-threads.ts`, `hooks/use-thread.ts`, `hooks/use-comments.ts`, `hooks/use-votes.ts` → page consumers  
**Why they're coupled**: Each page manually composes 4+ hooks. Duplicate fetch/transform logic. Violates `bundle-dynamic-imports` and potentially causes waterfalls.  
**Dependency category**: **Shared callback pattern** — hooks share similar API call patterns  
**Test impact**: Create boundary tests (mock API client) replacing individual hook tests that don't exist

**Architecture Pattern**: React Composition (Compound Hooks)

- **Current**: Each hook does independent fetch
- **Target**: Unified hook that composes sub-hooks, following `async-parallel` for independent operations

---

## Priority Order

| Priority | Candidate | Impact | Effort |
|----------|------------|--------|--------|
| 1 | Backend: Hexagonal Migration | High | Medium |
| 2 | Backend: Injectable Session Store | Medium | Low |
| 3 | Frontend: Unified Hooks | Medium | Medium |

---

## References

- Backend: Hexagonal Architecture pattern (Ports & Adapters)
- Frontend: Vercel React Best Practices (`vercel-react-best-practices` skill)
  - `async-parallel` - Use Promise.all() for independent operations
  - `bundle-dynamic-imports` - Use next/dynamic for heavy components
  - `rerender-memo` - Extract expensive work into memoized components