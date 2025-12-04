// internal/http/router.go
package httpserver

import (
	"net/http"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
	"github.com/whiterabbit0809/overengineered-calculator/internal/calculator"
	"github.com/whiterabbit0809/overengineered-calculator/internal/history"
)

func NewRouter(
	authHandler *auth.Handler,
	tokenService auth.TokenService,
	calcHandler *calculator.Handler,
	historyHandler *history.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	// Calculator (protected)
	mux.Handle("/api/v1/calc",
		Chain(http.HandlerFunc(calcHandler.Calculate), AuthMiddleware(tokenService)),
	)
	// History (protected)
	mux.Handle("/api/v1/history",
		Chain(http.HandlerFunc(historyHandler.GetHistory), AuthMiddleware(tokenService)),
	)

	// Static frontend
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	return mux
}
