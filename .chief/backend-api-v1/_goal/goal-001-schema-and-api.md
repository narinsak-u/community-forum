# Goal: Backend Schema & Full API Implementation

## Overview

Implement a complete backend API for the Community Forum with session-based authentication, full CRUD operations for threads/comments/users, tagging, voting, and view tracking. The backend connects to an existing React frontend that has UI pages for authentication, thread listing, thread creation, thread detail, and user profiles.

## Scope

This milestone covers the backend only — database schema, API handlers, session auth, and data-fetching endpoints. Frontend data wiring (TanStack Query hooks, form submissions) is out of scope.

## Constraints

- Go 1.24 + Fiber v2 + GORM + PostgreSQL
- No business logic in handlers — delegate to service layer
- All API responses use JSON; errors return `{ "error": "message" }`
- Follow AGENTS.md backend conventions: lowercase files, PascalCase exports, camelCase internals, UPPER_SNAKE_CASE constants
- No hardcoded connection strings — use environment variables
- Passwords hashed with bcrypt (cost 12)

## Data Model Decisions

| Entity | Key Decision |
|--------|------------|
| **User** | Has `username`, `email`, `password` (hashed), `avatar`, `bio`, `stacks` (tech stacks as JSON array), `role` (enum: user, admin) |
| **Thread** | Has `title`, `slug` (auto-generated from title), `content` (markdown text), `status` (enum: draft, published, archived), `viewCount` (denormalized), plus author, tags, comments, votes relationships |
| **Comment** | Nested via `parentId` (ParentID self-reference). Has `content`, `author`, `thread`, `parent`, `replies`, plus votes relationship |
| **Tag** | Has `name` (unique), `color`. Many-to-many with Thread via `thread_tags` |
| **Vote** | Tracks each user's vote on Thread or Comment. `value`: 1 (upvote), -1 (downvote). One vote per user per target |
| **Session** | GORM-managed session store for Fiber sessions middleware |

## Data Fetching Decisions

| Feature | Decision |
|---------|----------|
| **Featured thread** | Thread with highest vote score (`votes WHERE value=1` grouped by thread) within last 7 days |
| **Trending threads** | Top 3 threads by combined score (`upvotes - downvotes + reply_count * 0.5`) |
| **Thread listing** | Paginated (pageSize=5), ordered by `createdAt DESC`, filter by `status=published`, include author + tags + vote counts |
| **Thread detail** | Single thread by slug, include author + tags + top-level comments (with nested replies) + vote counts |
| **User threads** | All threads by userId, ordered by `createdAt DESC`, include vote counts |
| **Session auth** | Sessions stored in PostgreSQL via Fiber sessions middleware with secure cookies |

## Out of Scope

- Frontend TanStack Query hooks or API calls
- File/image upload (avatar upload)
- Email verification
- Password reset
- Rate limiting
- Admin panel
- Search functionality
- WebSocket/real-time updates