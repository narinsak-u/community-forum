package domain

import "time"

type Comment struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string
	ThreadID  uint
	AuthorID  uint
	Author    User
	ParentID  *uint
	Replies   []Comment
	Upvotes   int64
	Downvotes int64
}
