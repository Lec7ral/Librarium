// Package configs manages application configuration.
package configs

import (
	"log"
	"os"
	"strings"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort string
	AppHost    string // Hostname for the application (e.g., localhost:8080 or my-app.com)
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
	cfg.ServerPort = os.Getenv("PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = os.Getenv("SERVER_PORT")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	// --- App Host Configuration ---
	// Prioritize HOSTNAME, which is provided by the production environment (Domcloud).
	cfg.AppHost = os.Getenv("HOSTNAME")
	if cfg.AppHost == "" {
		// Fallback to APP_HOST for other environments.
		cfg.AppHost = os.Getenv("APP_HOST")
	}
	if cfg.AppHost == "" {
		// Default to localhost with the configured port for local development.
		port := cfg.ServerPort
		if strings.HasPrefix(port, ":") {
			port = port[1:]
		}
		cfg.AppHost = "localhost:" + port
	}

	// Ensure the ServerPort starts with a colon for ListenAndServe.
	if !strings.HasPrefix(cfg.ServerPort, ":") {
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
