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
	"github.com/Lec7ral/fullAPI/docs"
	"github.com/Lec7ral/fullAPI/internal/database"
	"github.com/Lec7ral/fullAPI/internal/handlers"
	"github.com/Lec7ral/fullAPI/internal/middleware"
	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Import godotenv
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	// The init function runs before main.
	// We load the .env file here. godotenv.Load() will NOT override existing environment variables.
	// This means that variables set by the hosting provider will always take precedence.
	// It will only load variables from .env if they are not already set in the environment.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environment variables.")
	}
}

// @title           Librarium API
// ... (resto de las anotaciones)
func main() {
	// By the time main starts, environment variables from .env (if present) are already loaded.

	// --- 1. SETUP ---
	cfg := configs.LoadConfig()

	// --- Dynamic Swagger Configuration ---
	docs.SwaggerInfo.Title = "Librarium API"
	docs.SwaggerInfo.Description = "This is the API for the Librarium application."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.PublicHost
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{cfg.PublicScheme}

	// ... (resto de main.go no cambia)
	db, err := database.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	bookRepo := repository.NewSQLiteBookRepository(db)
	userRepo := repository.NewSQLiteUserRepository(db)
	authorRepo := repository.NewSQLiteAuthorRepository(db)
	loanRepo := repository.NewSQLiteLoanRepository(db)
	env := &handlers.Env{
		BookRepo:   bookRepo,
		UserRepo:   userRepo,
		AuthorRepo: authorRepo,
		LoanRepo:   loanRepo,
		JWTSecret:  cfg.JWTSecret,
	}
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	authMw := middleware.AuthMiddleware(userRepo, cfg.JWTSecret)
	adminMw := middleware.RoleRequiredMiddleware("librarian")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/register", env.RegisterUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", env.LoginUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/authors", env.GetAuthorsHandler).Methods(http.MethodGet)
	router.HandleFunc("/authors/{id}", env.GetAuthorHandler).Methods(http.MethodGet)
	router.Handle("/authors", authMw(adminMw(http.HandlerFunc(env.CreateAuthorHandler)))).Methods(http.MethodPost)
	router.HandleFunc("/books", env.GetBooksHandler).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", env.GetBookHandler).Methods(http.MethodGet)
	router.Handle("/books", authMw(adminMw(http.HandlerFunc(env.CreateBookHandler)))).Methods(http.MethodPost)
	router.Handle("/books/{id}", authMw(adminMw(http.HandlerFunc(env.UpdateBookHandler)))).Methods(http.MethodPut)
	router.Handle("/books/{id}", authMw(adminMw(http.HandlerFunc(env.DeleteBookHandler)))).Methods(http.MethodDelete)
	router.Handle("/loans", authMw(http.HandlerFunc(env.CreateLoanHandler))).Methods(http.MethodPost)
	router.Handle("/loans/{id}", authMw(http.HandlerFunc(env.ReturnLoanHandler))).Methods(http.MethodDelete)
	router.Handle("/users/me/loans", authMw(http.HandlerFunc(env.GetMyLoansHandler))).Methods(http.MethodGet)
	router.Handle("/loans", authMw(adminMw(http.HandlerFunc(env.GetAllLoansHandler)))).Methods(http.MethodGet)
	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}
	go func() {
		log.Printf("Starting server on port %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting.")
}
