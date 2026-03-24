package domain

import "time"

// Note represents a persisted note record exposed to upper layers.
type Note struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
