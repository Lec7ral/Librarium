// Package repository provides a data abstraction layer.
// This file contains the implementation for loan data operations.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Lec7ral/fullAPI/internal/models"
)

// LoanFilter holds the criteria for searching loans.
type LoanFilter struct {
	Status *string // "active" or "returned"
}

// LoanRepository defines the interface for loan data operations.
type LoanRepository interface {
	CreateLoan(bookID, userID int64) error
	ReturnLoan(loanID int64) error
	GetActiveLoansByUserID(userID int64) ([]models.Loan, error)
	SearchLoans(filter LoanFilter) ([]models.Loan, error)
}

// sqliteLoanRepository is the concrete implementation for SQLite.
type sqliteLoanRepository struct {
	DB *sql.DB
}

// NewSQLiteLoanRepository creates a new repository instance.
func NewSQLiteLoanRepository(db *sql.DB) LoanRepository {
	return &sqliteLoanRepository{DB: db}
}

func (r *sqliteLoanRepository) CreateLoan(bookID, userID int64) error {
	tx, err := r.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStock int
	err = tx.QueryRow("SELECT stock FROM books WHERE id = ?", bookID).Scan(&currentStock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if currentStock <= 0 {
		return errors.New("no stock available")
	}

	_, err = tx.Exec("UPDATE books SET stock = stock - 1 WHERE id = ?", bookID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO loans (book_id, user_id, loan_date) VALUES (?, ?, ?)",
		bookID, userID, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *sqliteLoanRepository) ReturnLoan(loanID int64) error {
	tx, err := r.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var loan models.Loan
	var returnDate sql.NullTime
	err = tx.QueryRow("SELECT id, book_id, return_date FROM loans WHERE id = ?", loanID).Scan(&loan.ID, &loan.BookID, &returnDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if returnDate.Valid {
		return errors.New("book already returned")
	}

	_, err = tx.Exec("UPDATE loans SET return_date = ? WHERE id = ?", time.Now(), loanID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE books SET stock = stock + 1 WHERE id = ?", loan.BookID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *sqliteLoanRepository) GetActiveLoansByUserID(userID int64) ([]models.Loan, error) {
	q := `
		SELECT l.id, l.loan_date, l.book_id, b.title, b.isbn 
		FROM loans l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = ? AND l.return_date IS NULL
	`
	rows, err := r.DB.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		var book models.Book
		if err := rows.Scan(&loan.ID, &loan.LoanDate, &loan.BookID, &book.Title, &book.ISBN); err != nil {
			return nil, err
		}
		loan.Book = &book
		loans = append(loans, loan)
	}
	return loans, nil
}

// SearchLoans searches for loans with optional filters.
func (r *sqliteLoanRepository) SearchLoans(filter LoanFilter) ([]models.Loan, error) {
	query := `
		SELECT
			l.id, l.loan_date, l.return_date,
			b.id, b.title,
			u.id, u.username
		FROM loans l
		JOIN books b ON l.book_id = b.id
		JOIN users u ON l.user_id = u.id
	`
	var args []interface{}
	whereClause := " WHERE 1=1"

	if filter.Status != nil {
		if *filter.Status == "active" {
			whereClause += " AND l.return_date IS NULL"
		} else if *filter.Status == "returned" {
			whereClause += " AND l.return_date IS NOT NULL"
		}
	}

	rows, err := r.DB.Query(query+whereClause, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		var book models.Book
		var user models.User
		var returnDate sql.NullTime

		if err := rows.Scan(
			&loan.ID, &loan.LoanDate, &returnDate,
			&book.ID, &book.Title,
			&user.ID, &user.Username,
		); err != nil {
			return nil, err
		}

		// Correctly handle the nullable return date.
		if returnDate.Valid {
			loan.ReturnDate = &returnDate.Time
		}
		loan.Book = &book
		loan.User = &user
		loans = append(loans, loan)
	}
	return loans, nil
}
