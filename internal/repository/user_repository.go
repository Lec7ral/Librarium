// Package repository provides a data abstraction layer.
// This file contains the implementation for user data operations.
package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Lec7ral/fullAPI/internal/models"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user models.User, passwordHash string) error
	GetByUsername(username string) (*models.User, error)
}

// sqliteUserRepository is the concrete implementation for SQLite.
type sqliteUserRepository struct {
	DB *sql.DB
}

// NewSQLiteUserRepository creates a new repository instance.
func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{DB: db}
}

// Create inserts a new user into the database.
func (r *sqliteUserRepository) Create(user models.User, passwordHash string) error {
	stmt, err := r.DB.Prepare("INSERT INTO users (username, password_hash) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, passwordHash)
	if err != nil {
		// Check for a unique constraint violation and return the shared error.
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameExists
		}
		return err
	}
	return nil
}

// GetByUsername finds a user by their username.
func (r *sqliteUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password_hash FROM users WHERE username = ?"
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		// Use errors.Is to check for sql.ErrNoRows and return the shared ErrNotFound.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
