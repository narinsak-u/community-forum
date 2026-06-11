# Chat Improvement Plan: Performance, Security, and Reliability

This document outlines the tasks required to bring the WebSocket chat implementation up to production-grade standards following an idiomatic Go audit.

## 1. Objectives
- **Performance**: Eliminate broadcast bottlenecks by minimizing mutex lock contention.
- **Reliability**: Ensure resource cleanup and request cancellation by properly propagating contexts.
- **Security**: Prevent XSS by sanitizing user-generated content before broadcast.
- **Maintainability**: Clean up ignored errors and hardcoded constants.

## 2. Implementation Tasks

### Task 1: Optimize Broadcast Performance
- **Issue**: The current `broadcast` method holds `h.mu` while iterating and calling `WriteMessage`. A single slow client can block the entire chat server.
- **Fix**: Snapshot the list of connections under the lock, then iterate and write outside the lock.
- **File**: `backend/internal/handlers/chat.go`

### Task 2: Connection-Tied Contexts
- **Issue**: `HandleWebSocket` currently passes `context.Background()` to services. If a client disconnects during a slow DB query, the query continues to consume server resources.
- **Fix**: Use `context.WithCancel` tied to the WebSocket connection lifecycle.
- **File**: `backend/internal/handlers/chat.go`

### Task 3: Robust Error Handling
- **Issue**: Errors from `h.UserService.GetUserByID` are ignored with `_`.
- **Fix**: Properly check and log errors. If the user profile cannot be fetched, decide on a fallback (e.g., "Guest" or closing the connection if auth state is compromised).
- **File**: `backend/internal/handlers/chat.go`

### Task 4: Content Sanitization (XSS Prevention)
- **Issue**: Chat messages are broadcast exactly as sent, allowing for potential Script/HTML injection.
- **Fix**: Implement HTML sanitization in the `ChatService` layer.
- **Dependencies**: Consider adding `github.com/microcosm-cc/bluemonday`.
- **File**: `backend/internal/usecase/chat_service.go`

### Task 5: Consolidation of Constants
- **Issue**: Hardcoded limits (e.g., `15`) are scattered.
- **Fix**: Use the `DefaultChatLimit` constant consistently across the service and handlers.
- **Files**: `backend/internal/handlers/chat.go`, `backend/internal/usecase/chat_service.go`

## 3. Verification Plan
- **Concurrency Test**: Run a load test with 50+ clients where one client is intentionally "slow" (limited bandwidth) and verify other clients still receive messages instantly.
- **Security Test**: Attempt to send `<script>alert(1)</script>` and verify it is sanitized before being broadcast or saved.
- **Context Test**: Simulate a long-running DB query, disconnect the client, and verify the DB context is cancelled.
