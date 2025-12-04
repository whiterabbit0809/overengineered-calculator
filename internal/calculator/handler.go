package calculator

import (
	"encoding/json"
	"net/http"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

// Temporary protected endpoint.
// Will later be replaced by real calculator logic.
func SecretHelloHandler(w http.ResponseWriter, r *http.Request) {
	userID, email, ok := auth.UserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "no authenticated user in context",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from the protected calculator API!",
		"userId":  userID,
		"email":   email,
	})
}
