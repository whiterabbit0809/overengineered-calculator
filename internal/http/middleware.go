package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order: Chain(h, m1, m2) -> m2(m1(h))
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// AuthMiddleware builds a middleware that enforces a valid JWT in the Authorization header.
func AuthMiddleware(tokenService auth.TokenService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeUnauthorized(w, "invalid Authorization header format")
				return
			}

			tokenStr := parts[1]

			claims, err := tokenService.ParseToken(tokenStr)
			if err != nil {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			// Put user info into context so handlers can access it.
			ctx := auth.ContextWithUser(r.Context(), claims.UserID, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
