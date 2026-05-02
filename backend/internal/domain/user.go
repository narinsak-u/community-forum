package domain

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Email     string
	Password  string
	Avatar    string
	Bio       string
	Stacks    []string // Domain uses actual slice, not JSON string
	Role      string
}
