// Package web provides shared web-related utility functions.
package web

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON is a helper function to send consistent JSON responses.
// It is exported so it can be used by other packages (handlers, middleware).
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError is a helper function to send consistent JSON error responses.
// It is exported so it can be used by other packages.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
