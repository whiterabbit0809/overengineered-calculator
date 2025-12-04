package history

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

// NewHandler constructs a new History HTTP handler.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GetHistory handles GET /api/v1/history?limit=&offset=
//
// It:
//   - Reads the authenticated user (ID + email) from context.
//   - Parses limit/offset with defaults (20, 0).
//   - Asks the History service for that user's entries.
//   - Returns a JSON array where each item includes email.
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	// Pagination defaults
	limit := 20
	if v := q.Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if v := q.Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Get user identity from JWT/context
	userID, email, ok := auth.UserFromContext(ctx)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "unauthorized",
		})
		return
	}

	// Fetch history entries for this user
	entries, err := h.svc.List(ctx, userID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "could not load history",
		})
		return
	}

	// Map HistoryEntry -> historyResponseEntry and attach email
	resp := make([]historyResponseEntry, len(entries))
	for i, e := range entries {
		resp[i] = historyResponseEntry{
			ID:         e.ID,
			Expression: e.Expression,
			Result:     e.Result,
			CreatedAt:  e.CreatedAt.Format(time.RFC3339), // or another format if you prefer
			Email:      email,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
