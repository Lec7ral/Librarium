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
	// Check if we are running in the production environment (Domcloud/Passenger).
	if os.Getenv("IN_PASSENGER") == "1" {
		// If so, hardcode the known public URL. This is a pragmatic solution
		// when the environment cannot be configured with custom variables.
		cfg.PublicHost = "librarium.mnz.dom.my.id"
		cfg.PublicScheme = "https"
	} else {
		// Otherwise, default to local development settings.
		port := cfg.ServerPort
		if strings.HasPrefix(port, ":") {
			port = port[1:]
		}
		cfg.PublicHost = "localhost:" + port
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
