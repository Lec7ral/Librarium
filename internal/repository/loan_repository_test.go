// Package repository contains tests for the repository layer.
package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestCreateLoan_Success tests the successful transaction of creating a loan.
func TestCreateLoan_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteLoanRepository(db)
	bookID, userID := int64(1), int64(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT stock FROM books WHERE id = ?")).
		WithArgs(bookID).
		WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(5))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET stock = stock - 1 WHERE id = ?")).
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO loans (book_id, user_id, loan_date) VALUES (?, ?, ?)")).
		WithArgs(bookID, userID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.CreateLoan(bookID, userID)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateLoan_NoStock tests that the transaction is rolled back if there is no stock.
func TestCreateLoan_NoStock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteLoanRepository(db)
	bookID, userID := int64(1), int64(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT stock FROM books WHERE id = ?")).
		WithArgs(bookID).
		WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(0))
	mock.ExpectRollback()

	err = repo.CreateLoan(bookID, userID)

	if err == nil {
		t.Errorf("expected an error, but got nil")
	}
	if err.Error() != "no stock available" {
		t.Errorf("expected error 'no stock available', but got '%v'", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestReturnLoan_Success tests the successful transaction of returning a loan.
func TestReturnLoan_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteLoanRepository(db)
	loanID, bookID := int64(1), int64(5)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, book_id, return_date FROM loans WHERE id = ?")).
		WithArgs(loanID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "return_date"}).AddRow(loanID, bookID, nil))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE loans SET return_date = ? WHERE id = ?")).
		WithArgs(sqlmock.AnyArg(), loanID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET stock = stock + 1 WHERE id = ?")).
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.ReturnLoan(loanID)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestReturnLoan_AlreadyReturned tests the case where a loan has already been returned.
func TestReturnLoan_AlreadyReturned(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSQLiteLoanRepository(db)
	loanID, bookID := int64(1), int64(5)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, book_id, return_date FROM loans WHERE id = ?")).
		WithArgs(loanID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "book_id", "return_date"}).AddRow(loanID, bookID, time.Now()))
	mock.ExpectRollback()

	err = repo.ReturnLoan(loanID)

	if err == nil {
		t.Errorf("expected an error, but got nil")
	}
	if err.Error() != "book already returned" {
		t.Errorf("expected error 'book already returned', but got '%v'", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
