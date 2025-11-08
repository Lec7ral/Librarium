// Package handlers contains the HTTP handlers for the application.
// This file provides validation helper functions.
package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is a singleton instance of the validator.
var validate = validator.New()

// validationErrors processes a validator.ValidationErrors error and returns a map
// of human-readable error messages for each failed field.
func validationErrors(err error) map[string]string {
	// Create a map to hold the formatted error messages.
	errors := make(map[string]string)

	// Type assert the error to access the slice of validation errors.
	for _, err := range err.(validator.ValidationErrors) {
		// Get the field name in lowercase.
		field := strings.ToLower(err.Field())

		// Generate a descriptive error message based on the validation tag.
		switch err.Tag() {
		case "required":
			errors[field] = "This field is required."
		case "min":
			errors[field] = fmt.Sprintf("This field must be at least %s characters long.", err.Param())
		case "max":
			errors[field] = fmt.Sprintf("This field must be at most %s characters long.", err.Param())
		case "isbn":
			errors[field] = "This field must be a valid ISBN."
		case "datetime":
			errors[field] = fmt.Sprintf("This field must be in the format %s.", err.Param())
		default:
			errors[field] = "This field is invalid."
		}
	}
	return errors
}
