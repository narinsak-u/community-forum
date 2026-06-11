package domain

import "time"

// ChatMessage represents a message in the global chat room.
type ChatMessage struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `json:"author"`
}
