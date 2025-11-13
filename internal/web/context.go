// Package web provides shared web-related utility functions.
package web

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

// Defines the key for storing user information in the context.
const UserContextKey = ContextKey("user")
