# Real-time Chat (Discussions)

Date: 2026-06-09

## Overview

Add a real-time group chat page at `/discussions` where any authenticated user can send and receive messages in real time. The "Discussions" sidebar link currently pointing to `/thread/architectural-shift` is updated to `/discussions`.

## Architecture

### Backend: WebSocket on existing Fiber port

Uses `gofiber/contrib/websocket` (wraps `gorilla/websocket`) mounted on the existing Fiber app at `GET /ws/chat`. Same port, same JWT auth middleware. Follows the existing hexagonal architecture pattern.

**New dependencies:**
- `github.com/gofiber/contrib/websocket`
- `github.com/gorilla/websocket`

### Frontend: Custom hook + page

A `useChat` hook wraps the native `WebSocket` API. The discussions page at `app/(main)/discussions/page.tsx` uses the hook and renders the chat UI.

## Database Model

A `chat_messages` table stores all messages:

```
ChatMessage
  ID          uint      (PK)
  CreatedAt   time.Time
  AuthorID    uint      (FK → users.id)
  Author      User      (populated via GORM preload)
  Content     string    (text, not null)
```

Follows hexagonal layering: `domain.ChatMessage` → `ports.ChatRepository` → `db.GORM ChatRepository` → `usecase.ChatService` → `handlers.ChatHandler`.

## WebSocket Protocol

JSON messages over a single WebSocket connection.

### Server → Client

| Type | Payload | When |
|---|---|---|
| `init` | `{ messages: ChatMessage[], users: OnlineUser[] }` | On client connect |
| `message` | `{ id, content, author: { id, username, avatar }, created_at }` | New message from anyone |
| `user_joined` | `{ user: { id, username, avatar } }` | User connected |
| `user_left` | `{ user: { id, username, avatar } }` | User disconnected |
| `load_more` | `{ messages: ChatMessage[] }` | Response to pagination request |

### Client → Server

| Type | Payload | When |
|---|---|---|
| `message` | `{ content: string }` | User sends a message |
| `load_more` | `{ before: number }` | User scrolls to top (before = oldest visible message ID) |

### Pagination

- On connect, server sends the 15 most recent messages via `init`.
- When client scrolls to top, it sends `load_more`. Server queries `WHERE id < before ORDER BY id DESC LIMIT 15`, reverses to chronological order, and responds.
- No total count needed — infinite scroll backward.

### Online Presence

An in-memory `sync.Mutex`-protected map tracks connected user IDs. On connect, the user is added and `user_joined` is broadcast. On disconnect, the user is removed and `user_left` is broadcast.

## Frontend Components

### `useChat()` hook

```
returns { messages, onlineUsers, sendMessage, loadMore, isConnected, isLoading }
```

Lifecycle:
1. On mount (auth available) → connect to `ws://host/ws/chat`
2. `init` → set messages + onlineUsers, mark `isLoading = false`
3. `message` → append to messages
4. `user_joined` / `user_left` → update onlineUsers
5. `load_more` → prepend older messages
6. `sendMessage(content)` → send `{ type: "message", content }`
7. `loadMore()` → send `{ type: "load_more", before: oldestMessageId }`
8. On unmount → close connection

Auto-reconnect on drop (3 retries, exponential backoff: 1s, 2s, 4s).

### `/discussions` page

- Auth-gated via `useRequireAuth` (redirects to login).
- Chat message feed (center), auto-scrolls to bottom on new messages. Scroll-to-top triggers `loadMore`.
- Message input bar (bottom) with send button, disabled when not connected.
- Online users sidebar (right) with green dot indicator and username.

### Data types

```ts
interface ChatMessage {
  id: number;
  content: string;
  author: { id: number; username: string; avatar: string };
  created_at: string;
}

interface OnlineUser {
  id: number;
  username: string;
  avatar: string;
}
```

## Error Handling

- **Connection drop**: Auto-reconnect with backoff. `init` re-sent after reconnect.
- **Send failure**: No optimistic rendering — message only appears after server broadcast.
- **DB write failure**: Server logs error, message not broadcast. Acceptable loss for chat.
- **Unauthenticated WS upgrade**: Server sends 401 close code; client redirects to login.
- **Rate limiting**: Server enforces max 1 message per 500ms per user.

## Sidebar Update

Change `{ to: "/thread/architectural-shift", label: "Discussions" }` to `{ to: "/discussions", label: "Discussions" }` in `Sidebar.tsx`.

## Testing

- **Backend**: Unit tests for `ChatService` (send, load paginated) with mock repository.
- **Frontend**: Unit test for `useChat` with mock WebSocket. Render test for discussions page.

## Files Changed / Created

### Backend (new)
- `internal/domain/chat_message.go` — entity
- `internal/ports/chat.go` — repository + service interfaces
- `internal/adapters/db/chat_repository.go` — GORM implementation
- `internal/usecase/chat_service.go` — business logic
- `internal/handlers/chat.go` — WebSocket handler + connection manager
- `internal/models/models.go` — add ChatMessage model

### Backend (modified)
- `cmd/server/main.go` — wire dependencies, add WS route

### Frontend (new)
- `app/(main)/discussions/page.tsx` — chat page
- `hooks/use-chat.ts` — WebSocket hook

### Frontend (modified)
- `components/forge/Sidebar.tsx` — update Discussions link to `/discussions`
