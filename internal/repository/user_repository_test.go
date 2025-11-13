// Package repository contains tests for the repository layer.
package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lec7ral/fullAPI/internal/models"
)

// TestCreateUser_Success tests the successful creation of a user.
func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	user := models.User{Username: "testuser", Role: "member"}
	passwordHash := "hashed_password"

	// Expect the INSERT statement that now includes the 'role' column.
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)")).
		ExpectExec().
		WithArgs(user.Username, passwordHash, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user, passwordHash)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateUser_Duplicate tests the case where the username already exists.
func TestCreateUser_Duplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	user := models.User{Username: "testuser", Role: "member"}

	// Expect the INSERT statement that now includes the 'role' column.
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)")).
		ExpectExec().
		WithArgs(user.Username, "any_hash", user.Role).
		WillReturnError(errors.New("UNIQUE constraint failed: users.username"))

	err = repo.Create(user, "any_hash")

	if !errors.Is(err, ErrUsernameExists) {
		t.Errorf("expected error to be ErrUsernameExists, but got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetUserByUsername_Success tests the successful retrieval of a user.
func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	expectedUser := &models.User{ID: 1, Username: "testuser", PasswordHash: "hashed_password", Role: "member"}

	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "role"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.PasswordHash, expectedUser.Role)

	query := regexp.QuoteMeta("SELECT id, username, password_hash, role FROM users WHERE username = ?")
	mock.ExpectQuery(query).WithArgs("testuser").WillReturnRows(rows)

	user, err := repo.GetByUsername("testuser")

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if user.Username != expectedUser.Username {
		t.Errorf("expected username '%s', but got '%s'", expectedUser.Username, user.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetUserByUsername_NotFound tests the case where a user is not found.
func TestGetUserByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	query := regexp.QuoteMeta("SELECT id, username, password_hash, role FROM users WHERE username = ?")
	mock.ExpectQuery(query).WithArgs("nonexistent").WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByUsername("nonexistent")

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected error to be ErrNotFound, but got %v", err)
	}
	if user != nil {
		t.Errorf("expected a nil user, but got one")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
