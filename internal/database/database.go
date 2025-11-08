// Package database handles database initialization and schema creation.
package database

import (
	"database/sql"
	"log"

	// Import the sqlite3 driver with a blank identifier to register it.
	_ "github.com/mattn/go-sqlite3"
)

// InitDB connects to the database specified by the DSN and ensures the required tables exist.
// It returns the database connection pool (*sql.DB) or an error.
func InitDB(dsn string) (*sql.DB, error) {
	// Open a connection to the SQLite database.
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Prepare the SQL statement to create the 'authors' table if it doesn't exist.
	authorsTableStmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			bio TEXT
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = authorsTableStmt.Exec()
	if err != nil {
		return nil, err
	}

	// Prepare the SQL statement to create the 'users' table if it doesn't exist.
	usersTableStmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = usersTableStmt.Exec()
	if err != nil {
		return nil, err
	}

	// --- Book Table Migration ---
	// This simple migration drops the old table to recreate it with the new schema.
	// In a real production environment, a more sophisticated migration tool would be used.
	_, err = db.Exec("DROP TABLE IF EXISTS books")
	if err != nil {
		return nil, err
	}

	// Prepare the SQL statement to create the 'books' table with the new 'stock' column.
	booksTableStmt, err := db.Prepare(`
		CREATE TABLE books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			published_date TEXT NOT NULL,
			isbn TEXT UNIQUE NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			author_id INTEGER,
			FOREIGN KEY(author_id) REFERENCES authors(id)
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = booksTableStmt.Exec()
	if err != nil {
		return nil, err
	}

	log.Println("Database tables (re)created successfully.")
	return db, nil
}
