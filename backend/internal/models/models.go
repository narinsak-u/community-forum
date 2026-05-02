package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:255" json:"email"`
	Password  string         `gorm:"size:255" json:"-"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Stacks    string         `gorm:"type:jsonb" json:"-"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"`
	Sessions  []Session      `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Threads   []Thread       `gorm:"foreignKey:AuthorID" json:"threads,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
	Votes     []Vote         `gorm:"foreignKey:UserID" json:"votes,omitempty"`
}

type Session struct {
	ID        string    `gorm:"primaryKey;size:255" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      string    `gorm:"type:text" json:"data"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
}

type Thread struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string         `gorm:"size:255" json:"title"`
	Slug      string         `gorm:"uniqueIndex;size:255" json:"slug"`
	Content   string         `gorm:"type:text" json:"content"`
	Status    string         `gorm:"size:20;default:'draft'" json:"status"`
	ViewCount uint           `gorm:"default:0" json:"view_count"`
	AuthorID  uint           `gorm:"index" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Tags      []Tag          `gorm:"many2many:thread_tags;" json:"tags"`
	Comments  []Comment      `gorm:"foreignKey:ThreadID" json:"comments"`
	Votes     []Vote         `gorm:"foreignKey:ThreadID" json:"votes,omitempty"`
}

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Content   string         `gorm:"type:text" json:"content"`
	ThreadID  uint           `gorm:"index" json:"thread_id"`
	Thread    Thread         `gorm:"foreignKey:ThreadID" json:"thread"`
	AuthorID  uint           `gorm:"index" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"parent"`
	Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies"`
	Votes     []Vote         `gorm:"foreignKey:CommentID" json:"votes,omitempty"`
}

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"uniqueIndex;size:50" json:"name"`
	Color     string         `gorm:"size:20;default:'#6366f1'" json:"color"`
	Threads   []Thread       `gorm:"many2many:thread_tags;" json:"threads"`
}

type Vote struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID    uint           `gorm:"index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	ThreadID  *uint          `gorm:"index" json:"thread_id"`
	Thread    *Thread        `gorm:"foreignKey:ThreadID" json:"thread"`
	CommentID *uint          `gorm:"index" json:"comment_id"`
	Comment   *Comment       `gorm:"foreignKey:CommentID" json:"comment"`
	Value     int8           `gorm:"default:0" json:"value"`
}

func (t *Thread) Upvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("thread_id = ? AND value = 1", t.ID).Count(&count)
	return count
}

func (t *Thread) Downvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("thread_id = ? AND value = -1", t.ID).Count(&count)
	return count
}

func (t *Thread) RepliesCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&Comment{}).Where("thread_id = ? AND parent_id IS NULL", t.ID).Count(&count)
	return count
}

func (c *Comment) Upvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("comment_id = ? AND value = 1", c.ID).Count(&count)
	return count
}

func (c *Comment) Downvotes(db *gorm.DB) int64 {
	var count int64
	db.Model(&Vote{}).Where("comment_id = ? AND value = -1", c.ID).Count(&count)
	return count
}
