// Package models defines the data structures used throughout the application.
package models

import "time"

// Loan represents the structure of a loan in the library, including its nested book and user.
// It includes struct tags for JSON marshaling and validation.
type Loan struct {
	ID       int64     `json:"id"`
	BookID   int64     `json:"book_id"`
	UserID   int64     `json:"user_id"`
	LoanDate time.Time `json:"loan_date"`
	// ReturnDate is a pointer to time.Time to allow for null values from the database.
	// The omitempty tag ensures it's not included in the JSON response if it's null.
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Book       *Book      `json:"book,omitempty"`
	User       *User      `json:"user,omitempty"`
}
