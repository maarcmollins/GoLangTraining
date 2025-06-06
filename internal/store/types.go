package store

import "time"

// Item represents a single to-do entry.
type Item struct {
	ID          int       `json:"id"`          // unique integer ID
	Description string    `json:"description"` // the task text
	CreatedAt   time.Time `json:"created_at"`  // timestamp when added
}
