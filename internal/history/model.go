package history

import "time"

type HistoryEntry struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"userId"`
	Expression string    `json:"expression"`
	Result     float64   `json:"result"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Handler struct {
	svc Service
}

type historyResponseEntry struct {
	ID         int64   `json:"id"`
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	CreatedAt  string  `json:"createdAt"`
	Email      string  `json:"email"`
}
