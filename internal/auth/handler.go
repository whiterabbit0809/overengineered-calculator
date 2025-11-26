// internal/auth/handler.go
package auth

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{service: service}
}

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.service.SignUp(r.Context(), req.Email, req.Password); err != nil {
		// TODO: better error mapping
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(signUpResponse{
			Status:  "failed",
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(signUpResponse{
		Status: "ok",
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"` // "passed" or "failed"
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	passed, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{
		Status: "failed",
	}
	if passed {
		resp.Status = "passed"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
