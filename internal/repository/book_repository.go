// Package repository provides a data abstraction layer.
// This file contains the implementation for book data operations.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Lec7ral/fullAPI/internal/models"
)

// BookFilter holds the criteria for searching books.
type BookFilter struct {
	Title  *string
	Author *string
}

// BookRepository defines the interface for book data operations.
type BookRepository interface {
	Create(book models.Book) (int64, error)
	Update(id int64, book models.Book) error
	Delete(id int64) error
	GetByID(id int64) (*models.Book, error)
	Search(filter BookFilter, limit, offset int, sort, order string) ([]models.Book, int, error)
}

// sqliteBookRepository is the concrete implementation for SQLite.
type sqliteBookRepository struct {
	DB *sql.DB
}

// NewSQLiteBookRepository creates a new repository instance.
func NewSQLiteBookRepository(db *sql.DB) BookRepository {
	return &sqliteBookRepository{DB: db}
}

// Create now includes the stock field in the INSERT statement.
func (r *sqliteBookRepository) Create(book models.Book) (int64, error) {
	stmt, err := r.DB.Prepare("INSERT INTO books (title, published_date, isbn, stock, author_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(book.Title, book.PublishedDate, book.ISBN, book.Stock, book.AuthorID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update now includes the stock field in the UPDATE statement.
func (r *sqliteBookRepository) Update(id int64, book models.Book) error {
	stmt, err := r.DB.Prepare("UPDATE books SET title = ?, published_date = ?, isbn = ?, stock = ?, author_id = ? WHERE id = ?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(book.Title, book.PublishedDate, book.ISBN, book.Stock, book.AuthorID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete remains unchanged.
func (r *sqliteBookRepository) Delete(id int64) error {
	stmt, err := r.DB.Prepare("DELETE FROM books WHERE id = ?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// getBookWithAuthorSQL now selects the stock field.
const getBookWithAuthorSQL = `
	SELECT
		b.id, b.title, b.published_date, b.isbn, b.stock, b.author_id,
		a.id, a.name, a.bio
	FROM
		books b
	LEFT JOIN
		authors a ON b.author_id = a.id`

// scanBookAndAuthor now scans the stock field.
func scanBookAndAuthor(row *sql.Row) (*models.Book, error) {
	var book models.Book
	var author models.Author
	err := row.Scan(
		&book.ID, &book.Title, &book.PublishedDate, &book.ISBN, &book.Stock, &book.AuthorID,
		&author.ID, &author.Name, &author.Bio,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	book.Author = &author
	return &book, nil
}

// GetByID uses the updated helpers.
func (r *sqliteBookRepository) GetByID(id int64) (*models.Book, error) {
	query := getBookWithAuthorSQL + " WHERE b.id = ?"
	row := r.DB.QueryRow(query, id)
	return scanBookAndAuthor(row)
}

// Search is the fully updated method for searching, filtering, sorting, and paginating.
func (r *sqliteBookRepository) Search(filter BookFilter, limit, offset int, sort, order string) ([]models.Book, int, error) {
	// --- 1. Build the COUNT query ---
	countQuery := "SELECT COUNT(*) FROM books b LEFT JOIN authors a ON b.author_id = a.id WHERE 1=1"
	var countArgs []interface{}
	if filter.Title != nil {
		countQuery += " AND b.title LIKE ?"
		countArgs = append(countArgs, fmt.Sprintf("%%%s%%", *filter.Title))
	}
	if filter.Author != nil {
		countQuery += " AND a.name LIKE ?"
		countArgs = append(countArgs, fmt.Sprintf("%%%s%%", *filter.Author))
	}

	var totalRecords int
	err := r.DB.QueryRow(countQuery, countArgs...).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	// --- 2. Build the main data query ---
	query := getBookWithAuthorSQL + " WHERE 1=1"
	var args []interface{}
	if filter.Title != nil {
		query += " AND b.title LIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Title))
	}
	if filter.Author != nil {
		query += " AND a.name LIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Author))
	}

	allowedSortFields := map[string]bool{"title": true, "published_date": true, "author": true, "stock": true}
	if allowedSortFields[sort] {
		field := "b." + sort
		if sort == "author" {
			field = "a.name"
		}
		if strings.ToUpper(order) != "ASC" && strings.ToUpper(order) != "DESC" {
			order = "ASC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", field, order)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var author models.Author
		if err := rows.Scan(
			&book.ID, &book.Title, &book.PublishedDate, &book.ISBN, &book.Stock, &book.AuthorID,
			&author.ID, &author.Name, &author.Bio,
		); err != nil {
			return nil, 0, err
		}
		book.Author = &author
		books = append(books, book)
	}

	return books, totalRecords, nil
}
