// Package main is the entry point for the API application.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lec7ral/fullAPI/configs"
	"github.com/Lec7ral/fullAPI/internal/database"
	"github.com/Lec7ral/fullAPI/internal/handlers"
	"github.com/Lec7ral/fullAPI/internal/middleware"
	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	// --- 1. SETUP ---

	// Load application configuration from environment variables.
	cfg := configs.LoadConfig()

	// Initialize the database connection pool.
	db, err := database.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create an instance of the book repository.
	bookRepo := repository.NewSQLiteBookRepository(db)
	userRepo := repository.NewSQLiteUserRepository(db)
	authorRepo := repository.NewSQLiteAuthorRepository(db)

	// Create the environment struct, injecting all dependencies.
	env := &handlers.Env{
		BookRepo:   bookRepo,
		UserRepo:   userRepo,
		AuthorRepo: authorRepo,
		JWTSecret:  cfg.JWTSecret,
	}

	// --- 2. ROUTING ---

	// Create a new router instance.
	router := mux.NewRouter()

	// Apply global middleware that runs on every request.
	router.Use(middleware.LoggingMiddleware)

	// Define public routes that do not require authentication.
	router.HandleFunc("/register", env.RegisterUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", env.LoginUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/books", env.GetBooksHandler).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", env.GetBookHandler).Methods(http.MethodGet)
	router.HandleFunc("/authors", env.GetAuthorsHandler).Methods(http.MethodGet)
	router.HandleFunc("/authors/{id}", env.GetAuthorHandler).Methods(http.MethodGet)
	// Create a sub-router for routes that require authentication.
	protectedRoutes := router.NewRoute().Subrouter()

	// Use the middleware constructor to create and apply the auth middleware.
	// This injects the JWT secret key into the middleware.
	protectedRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // <-- CORRECCIÓN CLAVE

	// Define protected routes for writing data.
	protectedRoutes.HandleFunc("/books", env.CreateBookHandler).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/books/{id}", env.UpdateBookHandler).Methods(http.MethodPut)
	protectedRoutes.HandleFunc("/books/{id}", env.DeleteBookHandler).Methods(http.MethodDelete)
	protectedRoutes.HandleFunc("/authors", env.CreateAuthorHandler).Methods(http.MethodPost)

	// --- 3. GRACEFUL SHUTDOWN ---

	// Create a custom HTTP server for more control.
	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}

	// Start the server in a goroutine so it doesn't block.
	go func() {
		log.Printf("Starting server on port %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	// Create a channel to listen for OS shutdown signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a shutdown signal is received.
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout to allow ongoing requests to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting.")
}
