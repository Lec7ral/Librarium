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

// Create, Update, and Delete methods remain the same as before.
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

// GetByID now uses a 2-step query to avoid JOINs on a single-item lookup.
func (r *sqliteBookRepository) GetByID(id int64) (*models.Book, error) {
	// 1. Get the book
	var book models.Book
	query := "SELECT id, title, published_date, isbn, stock, author_id FROM books WHERE id = ?"
	err := r.DB.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.PublishedDate, &book.ISBN, &book.Stock, &book.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 2. Get the author if an author_id exists
	if book.AuthorID > 0 {
		var author models.Author
		authorQuery := "SELECT id, name, bio FROM authors WHERE id = ?"
		err = r.DB.QueryRow(authorQuery, book.AuthorID).Scan(&author.ID, &author.Name, &author.Bio)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		book.Author = &author
	}

	return &book, nil
}

// Search now uses a 2-query strategy to avoid the N+1 problem.
func (r *sqliteBookRepository) Search(filter BookFilter, limit, offset int, sort, order string) ([]models.Book, int, error) {
	// --- 1. Build the query for fetching book IDs that match the criteria ---
	var idArgs []interface{}
	idQuery := "SELECT b.id FROM books b"
	whereClause := " WHERE 1=1"

	if filter.Author != nil {
		idQuery += " JOIN authors a ON b.author_id = a.id"
		whereClause += " AND a.name LIKE ?"
		idArgs = append(idArgs, fmt.Sprintf("%%%s%%", *filter.Author))
	}
	if filter.Title != nil {
		whereClause += " AND b.title LIKE ?"
		idArgs = append(idArgs, fmt.Sprintf("%%%s%%", *filter.Title))
	}

	idQuery += whereClause

	// --- 2. Get the total count using the same filters ---
	countQuery := "SELECT COUNT(b.id) FROM books b"
	if strings.Contains(idQuery, "JOIN") {
		countQuery += " JOIN authors a ON b.author_id = a.id"
	}
	countQuery += whereClause

	var totalRecords int
	err := r.DB.QueryRow(countQuery, idArgs...).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	// --- 3. Apply sorting and pagination to the ID query ---
	allowedSortFields := map[string]bool{"title": true, "published_date": true, "stock": true}
	if allowedSortFields[sort] {
		if strings.ToUpper(order) != "ASC" && strings.ToUpper(order) != "DESC" {
			order = "ASC"
		}
		idQuery += fmt.Sprintf(" ORDER BY b.%s %s", sort, order)
	}
	idQuery += " LIMIT ? OFFSET ?"
	idArgs = append(idArgs, limit, offset)

	rows, err := r.DB.Query(idQuery, idArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookIDs []interface{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		bookIDs = append(bookIDs, id)
	}

	if len(bookIDs) == 0 {
		return []models.Book{}, totalRecords, nil
	}

	// --- 4. Fetch the full book and author data for the retrieved IDs ---
	mainQuery := getBookWithAuthorSQL + " WHERE b.id IN (?" + strings.Repeat(",?", len(bookIDs)-1) + ")"

	mainRows, err := r.DB.Query(mainQuery, bookIDs...)
	if err != nil {
		return nil, 0, err
	}
	defer mainRows.Close()

	booksMap := make(map[int64]*models.Book)
	for mainRows.Next() {
		var book models.Book
		var author models.Author
		if err := mainRows.Scan(
			&book.ID, &book.Title, &book.PublishedDate, &book.ISBN, &book.Stock, &book.AuthorID,
			&author.ID, &author.Name, &author.Bio,
		); err != nil {
			return nil, 0, err
		}
		book.Author = &author
		booksMap[book.ID] = &book
	}

	// Re-order the results to match the order of the bookIDs query.
	finalBooks := make([]models.Book, 0, len(bookIDs))
	for _, id := range bookIDs {
		if book, ok := booksMap[id.(int64)]; ok {
			finalBooks = append(finalBooks, *book)
		}
	}

	return finalBooks, totalRecords, nil
}

const getBookWithAuthorSQL = `
	SELECT
		b.id, b.title, b.published_date, b.isbn, b.stock, b.author_id,
		a.id, a.name, a.bio
	FROM
		books b
	LEFT JOIN
		authors a ON b.author_id = a.id`
