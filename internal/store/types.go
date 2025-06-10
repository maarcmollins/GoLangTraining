package store

import "time"

type ctxKey string

const TraceIDKey ctxKey = "traceID"

const (
	StatusNotStarted = "not started"
	StatusStarted    = "started"
	StatusCompleted  = "completed"
)

// Item represents a single to-do entry.
type Item struct {
	ID          int       `json:"id"`          // unique integer ID
	Description string    `json:"description"` // the task text
	CreatedAt   time.Time `json:"created_at"`  // timestamp when added
	Status      string    `json:"status"`      // status of the item
}
