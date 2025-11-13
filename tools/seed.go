//go:build ignore
// +build ignore

// This file is a standalone application to seed the database with a large volume of data.
// To run this seeder, use the command: go run seed.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting database super-seeder...")

	// --- 1. Connect to the Database ---
	db, err := sql.Open("sqlite3", "./library.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// --- 2. Check if Seeding is Necessary ---
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil && !strings.Contains(err.Error(), "no such table") {
		log.Fatalf("Failed to check book count: %v", err)
	}

	if count > 100 { // Only seed if the DB is empty or has very few books
		log.Println("Database already contains a large volume of data. Seeding is not required. Exiting.")
		return
	}

	log.Println("Database is empty or near-empty. Seeding with large volume of data...")

	// --- 3. Create Tables (Idempotent) ---
	createTablesSQL := `
	DROP TABLE IF EXISTS loans;
	DROP TABLE IF EXISTS books;
	DROP TABLE IF EXISTS authors;
	DROP TABLE IF EXISTS users;

	CREATE TABLE authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		bio TEXT
	);
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'member'
	);
	CREATE TABLE books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		published_date TEXT NOT NULL,
		isbn TEXT UNIQUE NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		author_id INTEGER,
		FOREIGN KEY(author_id) REFERENCES authors(id)
	);
	CREATE TABLE loans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		book_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		loan_date TEXT NOT NULL,
		return_date TEXT,
		FOREIGN KEY(book_id) REFERENCES books(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Use a transaction for performance when inserting many rows.
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Rollback on error

	// --- 4. Seed Authors ---
	log.Println("Seeding authors...")
	authorStmt, err := tx.Prepare("INSERT INTO authors (name, bio) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare author insert: %v", err)
	}
	defer authorStmt.Close()

	var authorIDs []int64
	for i := 1; i <= 50; i++ {
		name := fmt.Sprintf("Author %d", i)
		bio := fmt.Sprintf("This is the biography for Author %d.", i)
		result, err := authorStmt.Exec(name, bio)
		if err != nil {
			log.Fatalf("Failed to execute author insert: %v", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatalf("Failed to get author last insert ID: %v", err)
		}
		authorIDs = append(authorIDs, id)
	}
	log.Println("50 authors seeded successfully.")

	// --- 5. Seed Books ---
	log.Println("Seeding books...")
	bookStmt, err := tx.Prepare("INSERT INTO books (title, published_date, isbn, stock, author_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare book insert: %v", err)
	}
	defer bookStmt.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 500; i++ {
		title := fmt.Sprintf("Book Title %d", i)
		// Generate a random date in the last 50 years
		year := 1974 + r.Intn(50)
		month := 1 + r.Intn(12)
		day := 1 + r.Intn(28)
		publishedDate := fmt.Sprintf("%d-%02d-%02d", year, month, day)
		isbn := fmt.Sprintf("978-3-16-148410-%d", i)  // Fake but unique ISBN
		stock := r.Intn(20)                           // Random stock between 0 and 19
		authorID := authorIDs[r.Intn(len(authorIDs))] // Assign a random author

		_, err = bookStmt.Exec(title, publishedDate, isbn, stock, authorID)
		if err != nil {
			log.Fatalf("Failed to execute book insert: %v", err)
		}
	}
	log.Println("500 books seeded successfully.")

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Database seeding complete!")
}
