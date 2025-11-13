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
	// Be explicit about the role, even if the DB has a default.
	if user.Role == "" {
		user.Role = "member"
	}
	stmt, err := r.DB.Prepare("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, passwordHash, user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameExists
		}
		return err
	}
	return nil
}

// GetByUsername finds a user by their username and includes their role.
func (r *sqliteUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	// Updated query to select the 'role' column.
	query := "SELECT id, username, password_hash, role FROM users WHERE username = ?"
	// Updated scan to include the 'role' field.
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
