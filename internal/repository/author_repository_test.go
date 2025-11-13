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

// TestCreateAuthor_Success tests the successful creation of an author.
func TestCreateAuthor_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteAuthorRepository(db)
	authorToCreate := models.Author{
		Name: "George Orwell",
		Bio:  "English novelist, essayist, journalist and critic.",
	}

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO authors (name, bio) VALUES (?, ?)")).
		ExpectExec().
		WithArgs(authorToCreate.Name, authorToCreate.Bio).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdID, err := repo.Create(authorToCreate)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if createdID != 1 {
		t.Errorf("expected created ID to be 1, but got %d", createdID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAuthorByID_Success tests the successful retrieval of an author.
func TestGetAuthorByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteAuthorRepository(db)
	expectedAuthor := &models.Author{
		ID:   1,
		Name: "George Orwell",
		Bio:  "English novelist, essayist, journalist and critic.",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "bio"}).
		AddRow(expectedAuthor.ID, expectedAuthor.Name, expectedAuthor.Bio)

	query := regexp.QuoteMeta("SELECT id, name, bio FROM authors WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	author, err := repo.GetByID(1)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if author.Name != expectedAuthor.Name {
		t.Errorf("expected author name '%s', but got '%s'", expectedAuthor.Name, author.Name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAuthorByID_NotFound tests the case where an author is not found.
func TestGetAuthorByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteAuthorRepository(db)
	query := regexp.QuoteMeta("SELECT id, name, bio FROM authors WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(99).WillReturnError(sql.ErrNoRows)

	author, err := repo.GetByID(99)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected error to be ErrNotFound, but got %v", err)
	}
	if author != nil {
		t.Errorf("expected a nil author, but got one")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
