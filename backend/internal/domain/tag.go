package domain

import "time"

type Tag struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Color     string
}
