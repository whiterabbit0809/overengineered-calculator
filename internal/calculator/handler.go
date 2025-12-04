package calculator

import (
	"encoding/json"
	"net/http"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, _, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	res, err := h.svc.Calculate(r.Context(), userID, req)
	if err != nil {
		switch err {
		case ErrInvalidOperation:
			http.Error(w, `{"error":"invalid operation"}`, http.StatusBadRequest)
			return
		case ErrDivisionByZero:
			http.Error(w, `{"error":"division by zero"}`, http.StatusBadRequest)
			return
		default:
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
