package history

import "time"

type HistoryEntry struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"userId"`
	Expression string    `json:"expression"`
	Result     float64   `json:"result"`
	CreatedAt  time.Time `json:"createdAt"`
}
