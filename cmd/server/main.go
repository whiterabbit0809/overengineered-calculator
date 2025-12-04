// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/whiterabbit0809/overengineered-calculator/internal/auth"
	httpserver "github.com/whiterabbit0809/overengineered-calculator/internal/http"
	"github.com/whiterabbit0809/overengineered-calculator/internal/storage"
)

func main() {
	db, err := storage.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	userRepo := auth.NewPostgresUserRepository(db)
	hasher := auth.NewBcryptPasswordHasher()
	authService := auth.NewAuthService(userRepo, hasher)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	tokenService := auth.NewJWTTokenService(jwtSecret, "overengineered-calculator", 24*time.Hour)

	authHandler := auth.NewHandler(authService, tokenService)

	router := httpserver.NewRouter(authHandler, tokenService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
