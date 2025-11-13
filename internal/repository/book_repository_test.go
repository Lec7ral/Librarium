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

// TestGetByID_Success tests the successful retrieval of a book by its ID.
func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteBookRepository(db)

	expectedAuthor := &models.Author{ID: 1, Name: "Test Author"}
	expectedBook := &models.Book{ID: 1, Title: "Test Book", AuthorID: 1, Author: expectedAuthor}

	// Mock for the first query (get book)
	bookRows := sqlmock.NewRows([]string{"id", "title", "published_date", "isbn", "stock", "author_id"}).
		AddRow(expectedBook.ID, expectedBook.Title, "2023-01-01", "1234567890", 10, expectedBook.AuthorID)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, published_date, isbn, stock, author_id FROM books WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(bookRows)

	// Mock for the second query (get author)
	authorRows := sqlmock.NewRows([]string{"id", "name", "bio"}).
		AddRow(expectedAuthor.ID, expectedAuthor.Name, "A test bio")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, bio FROM authors WHERE id = ?")).
		WithArgs(expectedBook.AuthorID).
		WillReturnRows(authorRows)

	book, err := repo.GetByID(1)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if book == nil {
		t.Errorf("expected a book but got nil")
		return
	}
	if book.Title != expectedBook.Title {
		t.Errorf("expected title '%s' but got '%s'", expectedBook.Title, book.Title)
	}
	if book.Author == nil || book.Author.Name != expectedAuthor.Name {
		t.Errorf("expected author name '%s' but got '%s'", expectedAuthor.Name, book.Author.Name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetByID_NotFound tests the case where a book is not found.
func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteBookRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, published_date, isbn, stock, author_id FROM books WHERE id = ?")).
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)

	book, err := repo.GetByID(2)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected error to be ErrNotFound, but got %v", err)
	}
	if book != nil {
		t.Errorf("expected a nil book, but got one")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
