// Package repository provides a data abstraction layer.
// This file contains shared error variables and types for the repository layer.
package repository

import "errors"

// Shared error variables for the repository layer.
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUsernameExists = errors.New("username already exists")
)
