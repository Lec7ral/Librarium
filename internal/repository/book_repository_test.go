// Package repository provides a data abstraction layer.
// This file contains tests for the book repository.
package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lec7ral/fullAPI/internal/models"
)

// TestGetByID_Success tests the successful retrieval of a book by its ID.
func TestGetByID_Success(t *testing.T) {
	// --- 1. Setup ---
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteBookRepository(db)

	// --- 2. Define Expected Data and Mock Expectations ---

	// Define the nested author object we expect to receive.
	expectedAuthor := &models.Author{
		ID:   1,
		Name: "Alan Donovan & Brian Kernighan",
		Bio:  "Authors of the Go Programming Language book.",
	}

	// Define the book, now including the nested author struct.
	expectedBook := &models.Book{
		ID:            1,
		Title:         "The Go Programming Language",
		PublishedDate: "2015-10-26",
		ISBN:          "978-0134190440",
		AuthorID:      1,
		Author:        expectedAuthor,
	}

	// Define the columns that the JOIN query will return.
	rows := sqlmock.NewRows([]string{"id", "title", "published_date", "isbn", "author_id", "id", "name", "bio"}).
		AddRow(expectedBook.ID, expectedBook.Title, expectedBook.PublishedDate, expectedBook.ISBN, expectedBook.AuthorID, expectedAuthor.ID, expectedAuthor.Name, expectedAuthor.Bio)

	// The query now uses the JOIN defined in getBookWithAuthorSQL.
	query := regexp.QuoteMeta(getBookWithAuthorSQL + " WHERE b.id = ?")

	// Set up the mock expectation.
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	// --- 3. Execution ---
	book, err := repo.GetByID(1)

	// --- 4. Assertion ---
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if book == nil {
		t.Errorf("expected a book but got nil")
		return // Return to avoid panic on nil pointer
	}
	if book.Title != expectedBook.Title {
		t.Errorf("expected title '%s' but got '%s'", expectedBook.Title, book.Title)
	}
	// Assert on the nested author's name.
	if book.Author.Name != expectedBook.Author.Name {
		t.Errorf("expected author '%s' but got '%s'", expectedBook.Author.Name, book.Author.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetByID_NotFound tests the case where a book is not found.
func TestGetByID_NotFound(t *testing.T) {
	// --- 1. Setup ---
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteBookRepository(db)

	// --- 2. Define Mock Expectations ---
	// The query is the same JOIN query.
	query := regexp.QuoteMeta(getBookWithAuthorSQL + " WHERE b.id = ?")

	// We expect the query to return sql.ErrNoRows.
	mock.ExpectQuery(query).WithArgs(2).WillReturnError(sql.ErrNoRows)

	// --- 3. Execution ---
	book, err := repo.GetByID(2)

	// --- 4. Assertion ---
	if err == nil {
		t.Errorf("expected an error, but got nil")
	}

	// Check that the error is our specific ErrNotFound.
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected error to be ErrNotFound, but got %v", err)
	}

	if book != nil {
		t.Errorf("expected a nil book, but got a book with title '%s'", book.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
