package models

import "time"

// Loan represents the structure of a loan in the library, including its nested book and user.
// It includes struct tags for JSON marshaling and validation.
type Loan struct {
	ID         int64     `json:"id"`
	BookID     int64     `json:"book_id"`
	UserID     int64     `json:"user_id"`
	LoanDate   time.Time `json:"loan_date"`
	ReturnDate time.Time `json:"return_date"`
	Book       *Book     `json:"book,omitempty"`
	User       *User     `json:"user,omitempty"`
}
