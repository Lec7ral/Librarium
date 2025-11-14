//go:build ignore
// +build ignore

// This file is a standalone CLI tool to manage user roles.
// It is not part of the main API application and must be run manually.
//
// Usage Example:
// go run ./tools/manage_user.go --username="someuser" --role="librarian"

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Lec7ral/fullAPI/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// --- 1. Parse Command-Line Arguments (Flags) ---
	username := flag.String("username", "", "The username of the user to modify.")
	role := flag.String("role", "", "The new role to assign (e.g., 'librarian' or 'member').")
	flag.Parse()

	// Validate that the required flags were provided.
	if *username == "" || *role == "" {
		fmt.Println("Error: Both --username and --role flags are required.")
		flag.Usage() // Print the flag descriptions.
		os.Exit(1)   // Exit with an error code.
	}

	log.Printf("Attempting to set role for user '%s' to '%s'...", *username, *role)

	// --- 2. Connect to the Database ---
	db, err := sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatalf("FATAL: Failed to open database: %v", err)
	}
	defer db.Close()

	// --- 3. Use the Repository to Update the User ---
	userRepo := repository.NewSQLiteUserRepository(db)

	err = userRepo.UpdateUserRole(*username, *role)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Fatalf("FATAL: User '%s' not found.", *username)
		}
		log.Fatalf("FATAL: Failed to update user role: %v", err)
	}

	log.Printf("Success! User '%s' has been updated to role '%s'.", *username, *role)
}
