package domain

import "time"

type Vote struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	User      User
	ThreadID  *uint
	Thread    *Thread
	CommentID *uint
	Comment   *Comment
	Value     int8 // 1 for upvote, -1 for downvote, 0 for neutral
}
