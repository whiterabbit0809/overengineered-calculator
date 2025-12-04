// internal/history/model.go
package history

import "time"

type HistoryEntry struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"userId"`
	Expression string    `json:"expression"`
	Result     float64   `json:"result"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Handler wires HTTP requests to the History service.
type Handler struct {
	svc Service
}

// historyResponseEntry is the JSON shape returned by /api/v1/history.
// It is derived from HistoryEntry but uses a string for CreatedAt and
// includes the user's email (taken from the JWT/context).
type historyResponseEntry struct {
	ID         int64   `json:"id"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	CreatedAt  string  `json:"createdAt"`
	Email      string  `json:"email"`
}
