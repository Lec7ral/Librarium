// Package configs manages application configuration.
package configs

import (
	"log"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort string
	Database   struct{ DSN string }
	Redis      struct {
		Addr     string
		Password string
		DB       int
	}
	JWTSecret string
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() *Config {
	var cfg Config

	// --- Server Port Configuration ---
	// Most cloud providers (like Heroku, Domcloud, Render) set the PORT environment variable.
	// We will prioritize that variable for production compatibility.
	cfg.ServerPort = os.Getenv("PORT")
	if cfg.ServerPort == "" {
		// For local development or other setups, we can use a custom SERVER_PORT.
		cfg.ServerPort = os.Getenv("SERVER_PORT")
	}
	// If neither is set, use a default for local development.
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	// Ensure the port starts with a colon for the ListenAndServe function.
	if cfg.ServerPort[0] != ':' {
		cfg.ServerPort = ":" + cfg.ServerPort
	}

	// --- Database Configuration ---
	cfg.Database.DSN = os.Getenv("DB_DSN")
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = "./library.db"
	}

	// --- Redis Configuration ---
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "localhost:6379"
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	// --- JWT Configuration ---
	cfg.JWTSecret = os.Getenv("JWT_SECRET_KEY")
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "default_super_secret_key_for_dev_only"
	}

	log.Println("Configuration loaded")
	return &cfg
}
