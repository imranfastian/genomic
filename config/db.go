package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the shared database instance used throughout the app
var DB *sql.DB

// InitDB initializes the PostgreSQL connection
func InitDB() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "genomics")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("could not open DB: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("could not connect to DB: %w", err)
	}

	return nil
}

// CloseDB closes the active DB connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// getEnv returns a fallback value if the environment variable is not set
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
