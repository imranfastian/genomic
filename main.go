package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"genomic/config"
	"genomic/handlers"
	"genomic/middleware"
)

func main() {
	// Initialize PostgreSQL connection
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer config.CloseDB()

	mux := http.NewServeMux()

	// Public route
	mux.HandleFunc("/login", handlers.LoginHandler)

	// Protected route (uses JWT middleware)
	mux.HandleFunc("/genomes", middleware.JWTMiddleware(handlers.GenomesHandler))

	// Server config
	server := &http.Server{
		Addr:         ":8080",
		Handler:      logRequest(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Server is running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Middleware for logging HTTP requests
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
