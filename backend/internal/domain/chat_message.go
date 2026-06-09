package domain

import "time"

type ChatMessage struct {
	ID        uint
	CreatedAt time.Time
	Content   string
	AuthorID  uint
	Author    User
}
