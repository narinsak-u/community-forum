# Real-time Chat (Discussions) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a real-time group chat page at `/discussions` with WebSocket-based messaging and online presence.

**Architecture:** WebSocket mounted on the existing Fiber port via `gofiber/contrib/websocket`. Chat messages persisted in PostgreSQL via GORM. In-memory connection manager for online presence. Frontend uses native WebSocket API via a `useChat` hook.

**Tech Stack:** Go Fiber + gorilla/websocket, React 19 + TypeScript + Tailwind, PostgreSQL via GORM

---

### Task 1: Backend domain entity + model + ports

**Files:**
- Create: `backend/internal/domain/chat_message.go`
- Modify: `backend/internal/models/models.go`
- Create: `backend/internal/ports/chat.go`

- [ ] **Step 1: Create domain entity**

Write `backend/internal/domain/chat_message.go`:

```go
package domain

import "time"

type ChatMessage struct {
	ID        uint
	CreatedAt time.Time
	Content   string
	AuthorID  uint
	Author    User
}
```

- [ ] **Step 2: Add ChatMessage GORM model**

Add to `backend/internal/models/models.go` before the closing of the file (after Vote):

```go
type ChatMessage struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Content  string `gorm:"type:text;not null" json:"content"`
	AuthorID uint   `gorm:"index;not null" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID" json:"author"`
}
```

- [ ] **Step 3: Create ports interfaces**

Write `backend/internal/ports/chat.go`:

```go
package ports

import (
	"community-forum/backend/internal/domain"
	"context"
)

type ChatRepository interface {
	Save(ctx context.Context, msg *domain.ChatMessage) error
	GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}

type ChatService interface {
	SendMessage(ctx context.Context, authorID uint, content string) (*domain.ChatMessage, error)
	GetRecentMessages(ctx context.Context, limit int) ([]domain.ChatMessage, error)
	GetMessagesBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error)
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/domain/chat_message.go backend/internal/models/models.go backend/internal/ports/chat.go
git commit -m "feat: add chat domain entity, model, and ports"
```

---

### Task 2: Backend chat repository (GORM)

**Files:**
- Create: `backend/internal/adapters/db/chat_repository.go`
- Modify: `backend/internal/config/config.go` — add ChatMessage to AutoMigrate

- [ ] **Step 1: Create the GORM repository**

Write `backend/internal/adapters/db/chat_repository.go`:

```go
package db

import (
	"context"
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORM ChatRepository struct {
	db *gorm.DB
}

func NewGORM ChatRepository(db *gorm.DB) *GORM ChatRepository {
	return &GORM ChatRepository{db: db}
}

func (r *GORM ChatRepository) Save(ctx context.Context, msg *domain.ChatMessage) error {
	m := &models.ChatMessage{
		Content:  msg.Content,
		AuthorID: msg.AuthorID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	msg.CreatedAt = m.CreatedAt
	return nil
}

func (r *GORM ChatRepository) GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	var ms []models.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Order("id DESC").
		Limit(limit).
		Find(&ms).Error; err != nil {
		return nil, err
	}
	// Reverse to chronological order
	result := make([]domain.ChatMessage, 0, len(ms))
	for i := len(ms) - 1; i >= 0; i-- {
		result = append(result, *chatMessageFromModel(&ms[i]))
	}
	return result, nil
}

func (r *GORM ChatRepository) GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	var ms []models.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Where("id < ?", beforeID).
		Order("id DESC").
		Limit(limit).
		Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]domain.ChatMessage, 0, len(ms))
	for i := len(ms) - 1; i >= 0; i-- {
		result = append(result, *chatMessageFromModel(&ms[i]))
	}
	return result, nil
}

func chatMessageFromModel(m *models.ChatMessage) *domain.ChatMessage {
	return &domain.ChatMessage{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		Content:   m.Content,
		AuthorID:  m.AuthorID,
		Author: domain.User{
			ID:       m.Author.ID,
			Username: m.Author.Username,
			Avatar:   m.Author.Avatar,
		},
	}
}
```

- [ ] **Step 2: Add ChatMessage to AutoMigrate**

In `backend/internal/config/config.go`, add `&models.ChatMessage{}` to the `MigrateDB` function's `AutoMigrate` call:

```go
func MigrateDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Thread{},
		&models.Comment{},
		&models.Tag{},
		&models.Vote{},
		&models.ChatMessage{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/adapters/db/chat_repository.go backend/internal/config/config.go
git commit -m "feat: add GORM chat repository with auto-migration"
```

---

### Task 3: Backend chat use case (service)

**Files:**
- Create: `backend/internal/usecase/chat_service.go`
- Create: `backend/tests/usecase/chat_service_test.go` (first, TDD)
- Modify: `backend/tests/usecase/mock_repository_test.go` (add mock chat repo)

- [ ] **Step 1: Add mock chat repository to test helpers**

Add to `backend/tests/usecase/mock_repository_test.go` before `var errInternal`:

```go
type mockChatRepo struct {
	messages []domain.ChatMessage
	nextID   uint
	saveFn   func(msg *domain.ChatMessage) error
}

func newMockChatRepo() *mockChatRepo {
	return &mockChatRepo{
		messages: make([]domain.ChatMessage, 0),
		nextID:   1,
	}
}

func (m *mockChatRepo) Save(ctx context.Context, msg *domain.ChatMessage) error {
	if m.saveFn != nil {
		return m.saveFn(msg)
	}
	msg.ID = m.nextID
	m.nextID++
	msg.CreatedAt = time.Now()
	m.messages = append(m.messages, *msg)
	return nil
}

func (m *mockChatRepo) GetRecent(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	n := len(m.messages)
	if n == 0 {
		return []domain.ChatMessage{}, nil
	}
	start := n - limit
	if start < 0 {
		start = 0
	}
	result := make([]domain.ChatMessage, 0, n-start)
	for i := start; i < n; i++ {
		result = append(result, m.messages[i])
	}
	return result, nil
}

func (m *mockChatRepo) GetBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	result := make([]domain.ChatMessage, 0)
	for _, msg := range m.messages {
		if msg.ID < beforeID {
			result = append(result, msg)
		}
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}
```

Also add `"community-forum/backend/internal/ports"` import if not already there — and make sure the mock implements `ports.ChatRepository`. Let me check existing mocks — they use implicit interface satisfaction.

- [ ] **Step 2: Write the failing test**

Create `backend/tests/usecase/chat_service_test.go`:

```go
package usecase_test

import (
	"context"
	"testing"

	"community-forum/backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatService_SendMessage_Success(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	msg, err := svc.SendMessage(context.Background(), 1, "Hello world")
	require.NoError(t, err)
	assert.Equal(t, uint(1), msg.ID)
	assert.Equal(t, "Hello world", msg.Content)
	assert.Equal(t, uint(1), msg.AuthorID)
}

func TestChatService_GetRecentMessages_ReturnsLatest(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	_, _ = svc.SendMessage(context.Background(), 1, "first")
	_, _ = svc.SendMessage(context.Background(), 2, "second")
	_, _ = svc.SendMessage(context.Background(), 1, "third")

	messages, err := svc.GetRecentMessages(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, messages, 2)
	assert.Equal(t, "second", messages[0].Content)
	assert.Equal(t, "third", messages[1].Content)
}

func TestChatService_GetMessagesBefore_ReturnsOlder(t *testing.T) {
	mock := newMockChatRepo()
	svc := usecase.NewChatService(mock)

	m1, _ := svc.SendMessage(context.Background(), 1, "first")
	m2, _ := svc.SendMessage(context.Background(), 2, "second")
	_, _ = svc.SendMessage(context.Background(), 1, "third")

	messages, err := svc.GetMessagesBefore(context.Background(), m2.ID, 1)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	assert.Equal(t, m1.Content, messages[0].Content)
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./tests/usecase/ -run TestChatService -v`
Expected: FAIL — "undefined: usecase.NewChatService"

- [ ] **Step 4: Write the chat service implementation**

Write `backend/internal/usecase/chat_service.go`:

```go
package usecase

import (
	"context"
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/ports"
	"errors"
	"strings"
	"time"
)

type ChatService struct {
	repo ports.ChatRepository
}

func NewChatService(repo ports.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) SendMessage(ctx context.Context, authorID uint, content string) (*domain.ChatMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, errors.New("message content too long")
	}

	msg := &domain.ChatMessage{
		Content:  content,
		AuthorID: authorID,
	}
	if err := s.repo.Save(ctx, msg); err != nil {
		return nil, err
	}
	msg.CreatedAt = time.Now()
	return msg, nil
}

func (s *ChatService) GetRecentMessages(ctx context.Context, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = 15
	}
	return s.repo.GetRecent(ctx, limit)
}

func (s *ChatService) GetMessagesBefore(ctx context.Context, beforeID uint, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = 15
	}
	return s.repo.GetBefore(ctx, beforeID, limit)
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./tests/usecase/ -run TestChatService -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/internal/usecase/chat_service.go backend/tests/usecase/chat_service_test.go backend/tests/usecase/mock_repository_test.go
git commit -m "feat: add chat service with send and paginated load"
```

---

### Task 4: Add WebSocket dependencies

**Files:**
- Modify: `backend/go.mod` and `backend/go.sum` (via `go get`)

- [ ] **Step 1: Install WebSocket dependencies**

Run: `cd backend && go get github.com/gofiber/contrib/websocket`

This pulls in `gorilla/websocket` as a transitive dependency.

- [ ] **Step 2: Verify build**

Run: `cd backend && go build ./...`
Expected: SUCCESS

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "build: add gorilla/websocket and gofiber/contrib/websocket dependencies"
```

---

### Task 5: Backend WebSocket handler + connection manager

**Files:**
- Create: `backend/internal/handlers/chat.go`

- [ ] **Step 1: Add GetUserByID to UserService interface + impl**

Add to `backend/internal/ports/user.go` in the `UserService` interface:

```go
GetUserByID(ctx context.Context, id uint) (*domain.User, error)
```

Add to `backend/internal/usecase/user_service.go`:

```go
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
```

- [ ] **Step 2: Fix ChatRepository.Save to return populated author**

Update the `Save` method in `backend/internal/adapters/db/chat_repository.go`:

```go
func (r *GORM ChatRepository) Save(ctx context.Context, msg *domain.ChatMessage) error {
	m := &models.ChatMessage{
		Content:  msg.Content,
		AuthorID: msg.AuthorID,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Preload("Author").First(m, m.ID).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	msg.CreatedAt = m.CreatedAt
	msg.Author = domain.User{
		ID:       m.Author.ID,
		Username: m.Author.Username,
		Avatar:   m.Author.Avatar,
	}
	return nil
}
```

- [ ] **Step 3: Write the chat handler**

Write `backend/internal/handlers/chat.go`:

```go
package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/lib"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// --- JSON protocol types ---

type wsMessage struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Before  uint   `json:"before,omitempty"`
}

type wsOutgoing struct {
	Type     string              `json:"type"`
	Messages []domain.ChatMessage `json:"messages,omitempty"`
	Users    []onlineUser        `json:"users,omitempty"`
	User     *onlineUser         `json:"user,omitempty"`
	Message  *domain.ChatMessage `json:"message,omitempty"`
}

type onlineUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// --- Connection manager ---

type ChatHandler struct {
	Service       ports.ChatService
	SessionManager *middleware.SessionManager

	mu      sync.Mutex
	clients map[uint][]*websocket.Conn // userID -> connections
	users   map[uint]*onlineUser
}

func NewChatHandler(service ports.ChatService, sm *middleware.SessionManager) *ChatHandler {
	return &ChatHandler{
		Service:        service,
		SessionManager: sm,
		clients:        make(map[uint][]*websocket.Conn),
		users:          make(map[uint]*onlineUser),
	}
}

// UpgradeHandler upgrades HTTP to WebSocket. Applied as a Fiber handler before the WS middleware.
func (h *ChatHandler) UpgradeHandler(c *fiber.Ctx) error {
	tokenStr := c.Cookies("forge_token")
	if tokenStr == "" {
		auth := c.Get("Authorization")
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenStr = auth[7:]
		}
	}
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleWebSocket handles the WebSocket connection lifecycle.
func (h *ChatHandler) HandleWebSocket(c *websocket.Conn) {
	// Extract user from JWT cookie
	tokenStr := c.Cookies("forge_token")
	if tokenStr == "" {
		auth := c.Headers("Authorization")
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenStr = auth[7:]
		}
	}
	if tokenStr == "" {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4001, "unauthorized"))
		c.Close()
		return
	}

	claims, err := lib.VerifyJWT(tokenStr, h.SessionManager.Secret)
	if err != nil {
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4001, "unauthorized"))
		c.Close()
		return
	}

	userID := claims.UserID

	// Fetch user profile for username/avatar
	profile, _ := h.UserService.GetUserByID(c.Context(), userID)
	displayName := ""
	displayAvatar := ""
	if profile != nil {
		displayName = profile.Username
		displayAvatar = profile.Avatar
	}
	user := &onlineUser{ID: userID, Username: displayName, Avatar: displayAvatar}

	// Register connection
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], c)
	if _, exists := h.users[userID]; !exists {
		h.users[userID] = user
		h.broadcast(wsOutgoing{Type: "user_joined", User: user})
	}
	activeUsers := make([]onlineUser, 0, len(h.users))
	for _, u := range h.users {
		activeUsers = append(activeUsers, *u)
	}
	h.mu.Unlock()

	// Send init: recent messages + online users
	recent, _ := h.Service.GetRecentMessages(c.Context(), 15)
	initMsg := wsOutgoing{
		Type:     "init",
		Messages: recent,
		Users:    activeUsers,
	}
	if data, err := json.Marshal(initMsg); err == nil {
		_ = c.WriteMessage(websocket.TextMessage, data)
	}

	// Read loop
	lastMessageTime := time.Now()
	defer func() {
		h.mu.Lock()
		conns := h.clients[userID]
		for i, conn := range conns {
			if conn == c {
				h.clients[userID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
			delete(h.users, userID)
			h.broadcast(wsOutgoing{Type: "user_left", User: user})
		}
		h.mu.Unlock()
		c.Close()
	}()

	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			break
		}

		var msg wsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "message":
			// Rate limit: 1 msg per 500ms
			if time.Since(lastMessageTime) < 500*time.Millisecond {
				continue
			}
			lastMessageTime = time.Now()

			chatMsg, err := h.Service.SendMessage(c.Context(), userID, msg.Content)
			if err != nil {
				continue
			}
			h.broadcast(wsOutgoing{Type: "message", Message: chatMsg})

		case "load_more":
			older, err := h.Service.GetMessagesBefore(c.Context(), msg.Before, 15)
			if err != nil {
				continue
			}
			resp := wsOutgoing{Type: "load_more", Messages: older}
			if data, err := json.Marshal(resp); err == nil {
				_ = c.WriteMessage(websocket.TextMessage, data)
			}
		}
	}
}

func (h *ChatHandler) broadcast(msg wsOutgoing) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, conns := range h.clients {
		for _, conn := range conns {
			_ = conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
```

Note: The chat handler has a limitation — the online user list won't populate usernames until we wire in user fetching. Let me fix this by passing the UserService to ChatHandler as well. Let me update the design:

- [ ] **Step 2: Fix — ChatHandler needs UserService for username resolution**

Update the `ChatHandler` struct in `chat.go` to also accept `ports.UserService`, and when a user connects, fetch their profile to populate the `onlineUser`:

```go
type ChatHandler struct {
	Service        ports.ChatService
	UserService    ports.UserService
	SessionManager *middleware.SessionManager
	// ...
}

func NewChatHandler(service ports.ChatService, userService ports.UserService, sm *middleware.SessionManager) *ChatHandler {
```

In the connect handler, after verifying JWT, fetch the user profile:

```go
profile, err := h.UserService.GetUserProfile(c.Context(), "") // We need a GetByID method
```

Actually, let me simplify: we store the JWT claims which only have userID and role. We could add username to the JWT, or fetch it from DB on connect. Let me add a `GetByID` to the UserService. But that's scope creep.

Simpler approach: store the username in the JWT claims. Let me check the lib/jwt.go.

Actually, the simplest approach for presence is to just show user IDs for now, and fetch username from the message author data (which IS populated from DB). Let me simplify the presence to just track user IDs, and show usernames only when the user has sent a message.

Even simpler: fetch username from the chat handler by calling userService.GetByID. Let me add that method.

Let me keep the plan focused. I'll adjust: the chat handler fetches the user profile via UserService on connect.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handlers/chat.go
git commit -m "feat: add WebSocket chat handler with connection manager"
```

---

### Task 6: Wire chat into main.go

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Add chat handler initialization**

In `backend/cmd/server/main.go`, add after the other handler initializations (around line 91):

```go
chatRepo := db_adapter.NewGORM ChatRepository(db)
chatService := usecase.NewChatService(chatRepo)
chatHandler := handlers.NewChatHandler(chatService, userService, sessionManager)
```

- [ ] **Step 2: Add WebSocket route**

Add after the tag routes (before server start):

```go
app.Get("/ws/chat", chatHandler.UpgradeHandler, websocket.New(chatHandler.HandleWebSocket))
```

Make sure to import: `"github.com/gofiber/contrib/websocket"`

- [ ] **Step 3: Import the websocket package**

The import block should include `"github.com/gofiber/contrib/websocket"`:

```go
import (
	// ... existing imports ...
	"github.com/gofiber/contrib/websocket"
)
```

- [ ] **Step 4: Verify build**

Run: `cd backend && go build ./...`
Expected: SUCCESS

- [ ] **Step 5: Commit**

```bash
git add backend/cmd/server/main.go
git commit -m "feat: wire chat WebSocket handler and route in main.go"
```

---

### Task 7: Frontend useChat hook

**Files:**
- Create: `frontend/src/hooks/use-chat.ts`
- Create: `frontend/src/test/use-chat.test.ts`

- [ ] **Step 1: Write the hook**

Write `frontend/src/hooks/use-chat.ts`:

```ts
"use client";

import { useEffect, useRef, useState, useCallback } from "react";

export interface ChatMessage {
  id: number;
  content: string;
  author: { id: number; username: string; avatar: string };
  created_at: string;
}

export interface OnlineUser {
  id: number;
  username: string;
  avatar: string;
}

const WS_URL =
  (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace(/^http/, "ws") +
  "/api/v1/ws/chat";

interface UseChatReturn {
  messages: ChatMessage[];
  onlineUsers: OnlineUser[];
  sendMessage: (content: string) => void;
  loadMore: () => void;
  isConnected: boolean;
  isLoading: boolean;
}

export function useChat(): UseChatReturn {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<OnlineUser[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const wsRef = useRef<WebSocket | null>(null);
  const retriesRef = useRef(0);
  const maxRetries = 3;
  const oldestIdRef = useRef<number | null>(null);

  const connect = useCallback(() => {
    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      retriesRef.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        switch (data.type) {
          case "init":
            setMessages(data.messages || []);
            setOnlineUsers(data.users || []);
            setIsLoading(false);
            if (data.messages?.length) {
              oldestIdRef.current = data.messages[0].id;
            }
            break;
          case "message":
            setMessages((prev) => [...prev, data.message]);
            break;
          case "user_joined":
            setOnlineUsers((prev) => {
              if (prev.find((u) => u.id === data.user.id)) return prev;
              return [...prev, data.user];
            });
            break;
          case "user_left":
            setOnlineUsers((prev) =>
              prev.filter((u) => u.id !== data.user.id),
            );
            break;
          case "load_more":
            setMessages((prev) => [...data.messages, ...prev]);
            if (data.messages?.length) {
              oldestIdRef.current = data.messages[0].id;
            }
            break;
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
      if (retriesRef.current < maxRetries) {
        retriesRef.current++;
        const delay = Math.pow(2, retriesRef.current - 1) * 1000;
        setTimeout(connect, delay);
      }
    };

    ws.onerror = () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    connect();
    return () => {
      wsRef.current?.close();
    };
  }, [connect]);

  const sendMessage = useCallback((content: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: "message", content }));
    }
  }, []);

  const loadMore = useCallback(() => {
    if (oldestIdRef.current && wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({ type: "load_more", before: oldestIdRef.current }),
      );
    }
  }, []);

  return { messages, onlineUsers, sendMessage, loadMore, isConnected, isLoading };
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/hooks/use-chat.ts
git commit -m "feat: add useChat WebSocket hook"
```

---

### Task 8: Frontend discussions page

**Files:**
- Create: `frontend/src/app/(main)/discussions/page.tsx`

- [ ] **Step 1: Create the discussions page**

Write `frontend/src/app/(main)/discussions/page.tsx`:

```tsx
"use client";

import { useEffect, useRef, useState } from "react";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useChat, type ChatMessage, type OnlineUser } from "@/hooks/use-chat";
import { Skeleton } from "@/components/ui/skeleton";
import { Send } from "lucide-react";
import { timeAgo } from "@/lib/utils";

export default function DiscussionsPage() {
  const { requireAuth } = useRequireAuth();
  const { messages, onlineUsers, sendMessage, loadMore, isConnected, isLoading } = useChat();

  useEffect(() => {
    requireAuth({ redirect: "/discussions" });
  }, [requireAuth]);

  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);

  const handleSend = () => {
    const trimmed = input.trim();
    if (!trimmed) return;
    sendMessage(trimmed);
    setInput("");
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (autoScroll) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages, autoScroll]);

  // Detect scroll to top for pagination
  const handleScroll = () => {
    const el = messagesContainerRef.current;
    if (!el) return;
    setAutoScroll(el.scrollTop + el.clientHeight >= el.scrollHeight - 100);
    if (el.scrollTop === 0) {
      loadMore();
    }
  };

  return (
    <div className="flex h-[calc(100vh-4rem)]">
      {/* Chat area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <div className="border-b border-border/60 px-6 py-4">
          <h1 className="text-sm font-semibold text-foreground uppercase tracking-[0.12em]">
            # General Chat
          </h1>
        </div>

        {/* Messages */}
        <div
          ref={messagesContainerRef}
          onScroll={handleScroll}
          className="flex-1 overflow-y-auto px-6 py-4 space-y-4"
        >
          {isLoading ? (
            <div className="space-y-4 pt-4">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="flex gap-3">
                  <Skeleton className="h-8 w-8 rounded-full shrink-0" />
                  <div className="space-y-2 flex-1">
                    <Skeleton className="h-4 w-32" />
                    <Skeleton className="h-10 w-full max-w-md" />
                  </div>
                </div>
              ))}
            </div>
          ) : messages.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center pt-12">
              No messages yet. Start the conversation!
            </p>
          ) : (
            messages.map((msg: ChatMessage) => {
              const initials = (msg.author?.username || "??")
                .replace("@", "")
                .slice(0, 2)
                .toUpperCase();
              return (
                <div key={msg.id} className="flex gap-3">
                  <div className="h-8 w-8 rounded-full bg-gradient-signal grid place-items-center text-[10px] font-bold text-primary-foreground shrink-0">
                    {initials}
                  </div>
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-semibold text-foreground">
                        {msg.author?.username || "unknown"}
                      </span>
                      <span className="text-[10px] text-muted-foreground">
                        {msg.created_at ? timeAgo(msg.created_at) : ""}
                      </span>
                    </div>
                    <p className="text-sm text-foreground/90 leading-relaxed">
                      {msg.content}
                    </p>
                  </div>
                </div>
              );
            })
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* Input */}
        <div className="border-t border-border/60 px-6 py-4">
          <div className="flex gap-3">
            <input
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={isConnected ? "Type a message..." : "Connecting..."}
              disabled={!isConnected}
              className="flex-1 bg-secondary/40 border border-border rounded-sm px-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground/60 focus:outline-none focus:border-primary/40 disabled:opacity-50"
            />
            <button
              onClick={handleSend}
              disabled={!isConnected || !input.trim()}
              className="h-10 w-10 grid place-items-center bg-gradient-signal hover:opacity-90 text-primary-foreground rounded-sm disabled:opacity-40 disabled:cursor-not-allowed transition-all"
            >
              <Send className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Online users sidebar */}
      <div className="w-56 border-l border-border/60 flex flex-col shrink-0">
        <div className="border-b border-border/60 px-5 py-4">
          <div className="text-[10px] uppercase tracking-[0.18em] text-muted-foreground">
            Online <span className="text-foreground">— {onlineUsers.length}</span>
          </div>
        </div>
        <div className="flex-1 overflow-y-auto py-2">
          {onlineUsers.length === 0 ? (
            <p className="text-xs text-muted-foreground text-center pt-8">
              No users online
            </p>
          ) : (
            onlineUsers.map((u: OnlineUser) => {
              const initials = (u.username || "??")
                .replace("@", "")
                .slice(0, 2)
                .toUpperCase();
              return (
                <div
                  key={u.id}
                  className="flex items-center gap-3 px-5 py-2.5 hover:bg-secondary/30 transition-colors"
                >
                  <span className="h-2 w-2 rounded-full bg-success shrink-0" />
                  <div className="h-7 w-7 rounded-full bg-secondary grid place-items-center text-[9px] font-bold text-foreground shrink-0">
                    {initials}
                  </div>
                  <span className="text-xs text-foreground truncate">
                    {u.username || `user_${u.id}`}
                  </span>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/app/\(main\)/discussions/page.tsx
git commit -m "feat: add discussions chat page"
```

---

### Task 9: Update sidebar link

**Files:**
- Modify: `frontend/src/components/forge/Sidebar.tsx`

- [ ] **Step 1: Change the Discussions link**

In `Sidebar.tsx`, change:
```
{ to: "/thread/architectural-shift", label: "Discussions", icon: MessagesSquare },
```
to:
```
{ to: "/discussions", label: "Discussions", icon: MessagesSquare },
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/forge/Sidebar.tsx
git commit -m "feat: update sidebar Discussions link to /discussions"
```

---

### Task 10: Verify everything works

**Files:** None (verification only)

- [ ] **Step 1: Verify backend builds and tests pass**

Run: `cd backend && go build ./... && go test ./...`
Expected: BUILD SUCCESS, ALL TESTS PASS

- [ ] **Step 2: Verify frontend dev server starts**

Run: `cd frontend && timeout 15 bun run dev 2>&1 || true`
Expected: Dev server starts with no errors

- [ ] **Step 3: Final commit**

```bash
git add -A
git commit -m "chore: verify builds and tests pass"
```
