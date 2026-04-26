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
}

type Thread struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string         `gorm:"size:255" json:"title"`
	Slug      string         `gorm:"uniqueIndex;size:255" json:"slug"`
	Content   string         `gorm:"type:text" json:"content"`
	AuthorID  uint           `gorm:"index" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Tags      []Tag          `gorm:"many2many:thread_tags;" json:"tags"`
	Comments  []Comment      `gorm:"foreignKey:ThreadID" json:"comments"`
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
}

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"uniqueIndex;size:50" json:"name"`
	Color     string         `gorm:"size:20" json:"color"`
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
	Value     int8           `gorm:"default:0" json:"value"` // 1 for upvote, -1 for downvote
}
