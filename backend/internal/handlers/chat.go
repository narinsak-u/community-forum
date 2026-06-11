package handlers

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/lib"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"
	"community-forum/backend/internal/usecase"

	"github.com/gofiber/contrib/websocket"
)

type wsMessage struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Before  uint   `json:"before,omitempty"`
}

type wsOutgoing struct {
	Type     string               `json:"type"`
	Messages []domain.ChatMessage `json:"messages,omitempty"`
	Users    []onlineUser         `json:"users,omitempty"`
	User     *onlineUser          `json:"user,omitempty"`
	Message  *domain.ChatMessage  `json:"message,omitempty"`
}

type onlineUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ChatHandler struct {
	Service        ports.ChatService
	UserService    ports.UserService
	SessionManager *middleware.SessionManager

	mu      sync.Mutex
	clients map[uint][]*websocket.Conn
	users   map[uint]*onlineUser
}

func NewChatHandler(service ports.ChatService, userService ports.UserService, sm *middleware.SessionManager) *ChatHandler {
	return &ChatHandler{
		Service:        service,
		UserService:    userService,
		SessionManager: sm,
		clients:        make(map[uint][]*websocket.Conn),
		users:          make(map[uint]*onlineUser),
	}
}

func (h *ChatHandler) HandleWebSocket(c *websocket.Conn) {
	tokenStr := c.Cookies("forge_token")
	if tokenStr == "" {
		auth := c.Headers("Authorization")
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenStr = auth[7:]
		}
	}
	if tokenStr == "" {
		log.Printf("chat: missing token, closing")
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4001, "unauthorized"))
		c.Close()
		return
	}

	claims, err := lib.VerifyJWT(tokenStr, h.SessionManager.Secret)
	if err != nil {
		log.Printf("chat: invalid token: %v", err)
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4001, "unauthorized"))
		c.Close()
		return
	}

	userID := claims.UserID
	log.Printf("chat: user %d connected", userID)

	profile, _ := h.UserService.GetUserByID(context.Background(), userID)
	displayName := ""
	displayAvatar := ""
	if profile != nil {
		displayName = profile.Username
		displayAvatar = profile.Avatar
	}
	user := &onlineUser{ID: userID, Username: displayName, Avatar: displayAvatar}

	var isNewUser bool
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], c)
	if _, exists := h.users[userID]; !exists {
		h.users[userID] = user
		isNewUser = true
	}
	activeUsers := make([]onlineUser, 0, len(h.users))
	for _, u := range h.users {
		activeUsers = append(activeUsers, *u)
	}
	h.mu.Unlock()

	if isNewUser {
		h.broadcast(wsOutgoing{Type: "user_joined", User: user})
	}

	recent, err := h.Service.GetRecentMessages(context.Background(), usecase.DefaultChatLimit)
	if err != nil {
		log.Printf("chat: GetRecentMessages error: %v", err)
	}
	initMsg := wsOutgoing{
		Type:     "init",
		Messages: recent,
		Users:    activeUsers,
	}
	data, err := json.Marshal(initMsg)
	if err != nil {
		log.Printf("chat: marshal init error: %v", err)
	} else {
		if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("chat: write init error: %v", err)
		}
	}

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
		shouldBcast := len(h.clients[userID]) == 0
		if shouldBcast {
			delete(h.clients, userID)
			delete(h.users, userID)
		}
		h.mu.Unlock()

		if shouldBcast {
			h.broadcast(wsOutgoing{Type: "user_left", User: user})
		}
		c.Close()
		log.Printf("chat: user %d disconnected", userID)
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
			if time.Since(lastMessageTime) < 500*time.Millisecond {
				continue
			}
			lastMessageTime = time.Now()

			chatMsg, err := h.Service.SendMessage(context.Background(), userID, msg.Content)
			if err != nil {
				continue
			}
			h.broadcast(wsOutgoing{Type: "message", Message: chatMsg})

		case "load_more":
			older, err := h.Service.GetMessagesBefore(context.Background(), msg.Before, usecase.DefaultChatLimit)
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
		log.Printf("chat broadcast marshal error: %v", err)
		return
	}

	h.mu.Lock()
	allConns := make([]*websocket.Conn, 0)
	for _, conns := range h.clients {
		allConns = append(allConns, conns...)
	}
	h.mu.Unlock()

	for _, conn := range allConns {
		_ = conn.WriteMessage(websocket.TextMessage, data)
	}
}
