// Package repository provides a data abstraction layer.
// This file contains the implementation for author data operations.
package repository

import (
	"database/sql"
	"errors"

	"github.com/Lec7ral/fullAPI/internal/models"
)

// AuthorRepository defines the interface for author data operations.
type AuthorRepository interface {
	Create(author models.Author) (int64, error)
	GetAll() ([]models.Author, error)
	GetByID(id int64) (*models.Author, error)
}

// sqliteAuthorRepository is the concrete implementation for SQLite.
type sqliteAuthorRepository struct {
	DB *sql.DB
}

// NewSQLiteAuthorRepository creates a new repository instance.
func NewSQLiteAuthorRepository(db *sql.DB) AuthorRepository {
	return &sqliteAuthorRepository{DB: db}
}

func (r *sqliteAuthorRepository) Create(author models.Author) (int64, error) {
	stmt, err := r.DB.Prepare("INSERT INTO authors (name, bio) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(author.Name, author.Bio)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *sqliteAuthorRepository) GetAll() ([]models.Author, error) {
	query := "SELECT id, name, bio FROM authors"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []models.Author
	for rows.Next() {
		var author models.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Bio); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (r *sqliteAuthorRepository) GetByID(id int64) (*models.Author, error) {
	var author models.Author
	query := "SELECT id, name, bio FROM authors WHERE id = ?"
	err := r.DB.QueryRow(query, id).Scan(&author.ID, &author.Name, &author.Bio)
	if err != nil {
		// Use errors.Is to check for sql.ErrNoRows and return the shared ErrNotFound.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &author, nil
}
