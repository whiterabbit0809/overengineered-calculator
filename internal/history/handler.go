package history

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	// Parse pagination parameters with defaults: limit=20, offset=0
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

	// Build response with email and formatted time
	resp := make([]historyResponseEntry, len(entries))
	for i, e := range entries {
		resp[i] = historyResponseEntry{
			ID:         e.ID,
			Expression: e.Expression,
			Result:     e.Result,
			CreatedAt:  e.CreatedAt.Format(time.RFC3339),
			Email:      email,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
