// Package configs manages application configuration.
package configs

import (
	"log"
	"os"
	"strings"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort   string
	PublicHost   string // The public-facing hostname (e.g., my-app.com)
	PublicScheme string // The public-facing protocol (http or https)
	Database     struct{ DSN string }
	Redis        struct {
		Addr     string
		Password string
		DB       int
	}
	JWTSecret string
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
func LoadConfig() *Config {
	var cfg Config

	// --- Server Port (Internal) ---
	cfg.ServerPort = os.Getenv("PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if !strings.HasPrefix(cfg.ServerPort, ":") {
		cfg.ServerPort = ":" + cfg.ServerPort
	}

	// --- Public Host & Scheme (External) ---
	// This is what users (and Swagger) will see.
	cfg.PublicHost = os.Getenv("PUBLIC_HOST")
	if cfg.PublicHost == "" {
		// Default for local development
		cfg.PublicHost = "localhost" + cfg.ServerPort
	}
	cfg.PublicScheme = os.Getenv("PUBLIC_SCHEME")
	if cfg.PublicScheme == "" {
		// Default for local development
		cfg.PublicScheme = "http"
	}

	// --- Database, Redis, JWT Configurations ---
	cfg.Database.DSN = os.Getenv("DB_DSN")
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = "./library.db"
	}
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "localhost:6379"
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.JWTSecret = os.Getenv("JWT_SECRET_KEY")
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "default_super_secret_key_for_dev_only"
	}

	log.Println("Configuration loaded")
	return &cfg
}
