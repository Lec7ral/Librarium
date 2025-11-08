// Package models defines the data structures used throughout the application.
package models

// User represents a user account in the system.
// It includes struct tags for JSON marshaling and validation.
type User struct {
	// ID is the unique identifier for the user.
	ID int64 `json:"id"`
	// Username is the unique name for the user account.
	Username string `json:"username" validate:"required,min=3,max=50"`
	// PasswordHash is the hashed version of the user's password.
	// The json:"-" tag ensures this field is never exposed in API responses.
	PasswordHash string `json:"-"`
}
