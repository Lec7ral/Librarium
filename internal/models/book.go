// Package models defines the data structures used throughout the application.
package models

// Book represents the structure of a book in the library, including its nested author.
// It includes struct tags for JSON marshaling and validation.
type Book struct {
	// ID is the unique identifier for the book.
	ID int64 `json:"id"`
	// Title is the title of the book.
	Title string `json:"title" validate:"required,min=2,max=100"`
	// PublishedDate is the date the book was published, in YYYY-MM-DD format.
	PublishedDate string `json:"published_date" validate:"required,datetime=2006-01-02"`
	// ISBN is the International Standard Book Number.
	ISBN string `json:"isbn" validate:"required,isbn"`
	// Stock is the number of available copies of the book.
	Stock int `json:"stock" validate:"gte=0"` // gte=0 means "greater than or equal to 0"

	// AuthorID is used for data input when creating/updating a book.
	// It links the book to an author in the authors table.
	AuthorID int64 `json:"author_id" validate:"required"`

	// Author is a pointer to the Author model for nesting author information in responses.
	// The `omitempty` tag prevents it from being included in the JSON if it's nil.
	Author *Author `json:"author,omitempty"`
}
