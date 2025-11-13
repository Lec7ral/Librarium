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
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	// JWTSecret is the secret key for signing JWT tokens.
	JWTSecret string
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
// It provides default values for development if environment variables are not set.
func LoadConfig() *Config {
	var cfg Config

	// Load server port from environment or use default.
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = ":8080"
	}
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = ":6379"
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	// Load database DSN from environment or use default.
	cfg.Database.DSN = os.Getenv("DB_DSN")
	if cfg.Database.DSN == "" {
		cfg.Database.DSN = "./library.db"
	}

	// Load JWT secret key from environment or use a default (unsafe for production).
	cfg.JWTSecret = os.Getenv("JWT_SECRET_KEY")
	if cfg.JWTSecret == "" {
		// THIS DEFAULT VALUE IS INSECURE AND FOR DEVELOPMENT ONLY.
		// A real production environment should fail to start if this key is not provided.
		cfg.JWTSecret = "default_super_secret_key_for_dev_only"
	}

	log.Println("Configuration loaded")
	return &cfg
}
