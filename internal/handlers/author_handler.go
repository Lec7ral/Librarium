// Package handlers contains the HTTP handlers for the application.
// This file contains the handlers for author-related CRUD operations.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Lec7ral/fullAPI/internal/models"
	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/Lec7ral/fullAPI/internal/web" // <-- Import the new web package
	"github.com/gorilla/mux"
)

// CreateAuthorHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var newAuthor models.Author
	if err := json.NewDecoder(r.Body).Decode(&newAuthor); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(newAuthor); err != nil {
		errors := validationErrors(err)
		web.RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errors})
		return
	}

	id, err := e.AuthorRepo.Create(newAuthor)
	if err != nil {
		log.Printf("Handler error creating author: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Failed to create author")
		return
	}

	createdAuthor, err := e.AuthorRepo.GetByID(id)
	if err != nil {
		log.Printf("Handler error fetching created author: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, createdAuthor)
}

// GetAuthorsHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) GetAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := e.AuthorRepo.GetAll()
	if err != nil {
		log.Printf("Handler error getting all authors: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if authors == nil {
		authors = []models.Author{}
	}

	web.RespondWithJSON(w, http.StatusOK, authors)
}

// GetAuthorHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) GetAuthorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	author, err := e.AuthorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Author not found")
		} else {
			log.Printf("Handler error getting author by ID: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	web.RespondWithJSON(w, http.StatusOK, author)
}
