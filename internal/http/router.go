// internal/http/router.go
package httpserver

import (
	"net/http"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

func NewRouter(authHandler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	// --- API routes ---
	mux.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	// --- Frontend route ---
	// Serve index.html explicitly on "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "web/index.html")
	})

	return mux
}
