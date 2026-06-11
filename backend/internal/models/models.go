package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a forum member.
// In GORM, a struct becomes a database table. Each field becomes a column.
type User struct {
	// gorm:"primarykey" tells the database this is the unique identifier for each row.
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"` // Automatically managed by GORM
	UpdatedAt time.Time `json:"updated_at"` // Automatically managed by GORM
	// DeletedAt enables "Soft Delete". Rows aren't actually deleted, just marked as deleted.
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username string `gorm:"uniqueIndex;size:50" json:"username"` // uniqueIndex ensures no duplicate usernames
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
	Password string `gorm:"size:255" json:"-"` // json:"-" means this field is never sent in API responses
	Avatar   string `gorm:"size:255" json:"avatar"`
	Bio      string `gorm:"type:text" json:"bio"`
	Stacks   string `gorm:"type:jsonb" json:"-"` // Stored as JSON in PostgreSQL
	Role     string `gorm:"size:20;default:'user'" json:"role"`

	// Relationships: These define how tables link to each other.
	// foreignKey tells GORM which column in the other table points back to this User.
	Sessions []Session `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Threads  []Thread  `gorm:"foreignKey:AuthorID" json:"threads,omitempty"`
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
	Votes    []Vote    `gorm:"foreignKey:UserID" json:"votes,omitempty"`
}

// Session represents a logged-in state for a user.
type Session struct {
	ID        string    `gorm:"primaryKey;size:255" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      string    `gorm:"type:text" json:"data"` // Stores serialized session data
	UserID    uint      `gorm:"index" json:"user_id"`  // Links this session to a specific user
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
}

// Thread represents a discussion topic created by a user.
type Thread struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Title     string `gorm:"size:255" json:"title"`
	Slug      string `gorm:"uniqueIndex;size:255" json:"slug"` // URL-friendly version of the title
	Content   string `gorm:"type:text" json:"content"`
	Status    string `gorm:"size:20;default:'draft'" json:"status"` // e.g., 'published' or 'draft'
	ViewCount uint   `gorm:"default:0" json:"view_count"`

	AuthorID uint `gorm:"index" json:"author_id"`
	Author   User `gorm:"foreignKey:AuthorID" json:"author"`

	// many2many creates a "join table" (thread_tags) to link multiple threads to multiple tags.
	Tags     []Tag     `gorm:"many2many:thread_tags;" json:"tags"`
	Comments []Comment `gorm:"foreignKey:ThreadID" json:"comments"`
	Votes    []Vote    `gorm:"foreignKey:ThreadID" json:"votes,omitempty"`

	// Computed fields — populated via subquery SELECTs, not stored in DB.
	Upvotes      int64 `gorm:"<-:false" json:"upvotes"`
	Downvotes    int64 `gorm:"<-:false" json:"downvotes"`
	RepliesCount int64 `gorm:"<-:false" json:"replies_count"`
}

// Comment represents a reply within a thread.
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Content  string `gorm:"type:text" json:"content"`
	ThreadID uint   `gorm:"index" json:"thread_id"`
	Thread   Thread `gorm:"foreignKey:ThreadID" json:"thread"`
	AuthorID uint   `gorm:"index" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID" json:"author"`

	// Self-referencing relationship: A comment can have a parent comment (nested replies).
	ParentID *uint     `gorm:"index" json:"parent_id"`
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies"`

	Votes []Vote `gorm:"foreignKey:CommentID" json:"votes,omitempty"`
}

// Tag represents a category label (e.g., "Go", "React").
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name    string   `gorm:"uniqueIndex;size:50" json:"name"`
	Color   string   `gorm:"size:20;default:'#6366f1'" json:"color"`
	Threads []Thread `gorm:"many2many:thread_tags;" json:"threads"`
}

// Vote represents a user's upvote or downvote on a thread or comment.
type Vote struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint `gorm:"uniqueIndex:idx_vote_thread;uniqueIndex:idx_vote_comment" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	// Pointers (*uint) allow these fields to be NULL in the database.
	// A vote belongs to EITHER a thread OR a comment.
	// Unique indexes on nullable columns allow multiple NULLs (one vote per user per thread/comment).
	ThreadID  *uint    `gorm:"uniqueIndex:idx_vote_thread" json:"thread_id"`
	Thread    *Thread  `gorm:"foreignKey:ThreadID" json:"thread"`
	CommentID *uint    `gorm:"uniqueIndex:idx_vote_comment" json:"comment_id"`
	Comment   *Comment `gorm:"foreignKey:CommentID" json:"comment"`

	Value int8 `gorm:"default:0" json:"value"` // 1 for upvote, -1 for downvote, 0 for neutral
}

type ChatMessage struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Content  string `gorm:"type:text;not null" json:"content"`
	AuthorID uint   `gorm:"index;not null" json:"author_id"`
	Author   User   `gorm:"foreignKey:AuthorID" json:"author"`
}

// Upvotes counts all positive votes for a comment.
func (c *Comment) Upvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("comment_id = ? AND value = 1", c.ID).Count(&count)
	return count
}

// Downvotes counts all negative votes for a comment.
func (c *Comment) Downvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("comment_id = ? AND value = -1", c.ID).Count(&count)
	return count
}
