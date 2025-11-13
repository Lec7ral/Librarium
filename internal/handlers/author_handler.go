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
	"github.com/Lec7ral/fullAPI/internal/web"
	"github.com/gorilla/mux"
)

// @Summary      Create a new author
// @Description  Adds a new author to the collection. Requires librarian role.
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        author  body      models.Author  true  "Author object to be created"
// @Success      201     {object}  models.Author
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Security     BearerAuth
// @Router       /authors [post]
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

// @Summary      List authors
// @Description  Get a list of all authors.
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Author
// @Failure      500  {object}  map[string]string
// @Router       /authors [get]
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

// @Summary      Get an author by ID
// @Description  Retrieves the details of a single author by their unique ID.
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Author ID"
// @Success      200  {object}  models.Author
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /authors/{id} [get]
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
