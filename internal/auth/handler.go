// internal/auth/handler.go
package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

func NewHandler(service AuthService, tokenService TokenService) *Handler {
	return &Handler{
		service:      service,
		tokenService: tokenService,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(signUpResponse{
			Status:  "failed",
			Message: "method not allowed",
		})
		return
	}

	var req signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(signUpResponse{
			Status:  "failed",
			Message: "invalid request body",
		})
		return
	}

	if err := h.service.SignUp(r.Context(), req.Email, req.Password); err != nil {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case errors.Is(err, ErrInvalidEmailFormat):
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(signUpResponse{
				Status:  "failed",
				Field:   "email",
				Message: "invalid email format",
			})
		case errors.Is(err, ErrPasswordTooShort):
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(signUpResponse{
				Status:  "failed",
				Field:   "password",
				Message: "password must be at least 8 characters",
			})
		case errors.Is(err, ErrPasswordTooWeak):
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(signUpResponse{
				Status:  "failed",
				Field:   "password",
				Message: "password must contain at least one letter and one digit",
			})
		case errors.Is(err, ErrEmailAlreadyExists):
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(signUpResponse{
				Status:  "failed",
				Field:   "email",
				Message: "email already exists",
			})
		default:
			// Unexpected internal error
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(signUpResponse{
				Status:  "failed",
				Message: "internal error",
			})
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(signUpResponse{
		Status: "ok",
	})
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

	ok, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{Status: "failed"}

	if ok {
		// Load the user (we know they exist & password was correct)
		user, err := h.service.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Generate JWT token
		token, err := h.tokenService.GenerateToken(user)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp.Status = "passed"
		resp.Token = token
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
