// Package configs manages application configuration.
package configs

import (
	"log"
	"os"
)

// Config holds all configuration for the application.
// Values are loaded from environment variables.
type Config struct {
	// ServerPort defines the port for the HTTP server.
	ServerPort string
	// Database holds database-specific configuration.
	Database struct {
		// DSN is the Data Source Name for the database connection.
		DSN string
	}
	// Redis holds Redis-specific configuration.
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	// JWTSecret is the secret key for signing JWT tokens.
	JWTSecret string
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() *Config {
	var cfg Config

	// --- Server Port Configuration ---
	// Check for SERVER_PORT first, then fall back to PORT (common in production environments).
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = os.Getenv("PORT")
	}
	// If neither is set, use a default for local development.
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	// Ensure the port starts with a colon.
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
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD") // Default is no password
	// You can also load cfg.Redis.DB from an env var if needed.

	// --- JWT Configuration ---
	cfg.JWTSecret = os.Getenv("JWT_SECRET_KEY")
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "default_super_secret_key_for_dev_only"
	}

	log.Println("Configuration loaded")
	return &cfg
}
