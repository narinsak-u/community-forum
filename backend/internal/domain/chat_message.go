package domain

import "time"

// ChatMessage represents a message in the global chat room.
type ChatMessage struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string
	AuthorID  uint
	Author    User
}
