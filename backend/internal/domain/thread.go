package domain

import "time"

type Thread struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Title            string
	Slug             string
	Content          string
	Status           string
	ViewCount        uint
	AuthorID         uint
	Author           User
	Tags             []Tag
	Comments         []Comment
	Upvotes          int64
	Downvotes        int64
	RepliesCount     int64
	RecentCommenters []User
}
