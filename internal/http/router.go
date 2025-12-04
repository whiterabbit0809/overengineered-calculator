// internal/http/router.go
package httpserver

import (
	"net/http"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
	"github.com/whiterabbit0809/overengineered-calculator/internal/calculator"
)

func NewRouter(authHandler *auth.Handler, tokenService auth.TokenService) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	// Protected route: requires valid JWT
	mux.Handle("/api/v1/secret-hello",
		Chain(
			http.HandlerFunc(calculator.SecretHelloHandler),
			AuthMiddleware(tokenService),
		))

	// Static frontend
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	return mux
}
