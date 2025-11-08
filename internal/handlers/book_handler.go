// Package handlers contains the HTTP handlers for the application.
// This file contains the handlers for book-related CRUD operations.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/Lec7ral/fullAPI/internal/models"
	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/Lec7ral/fullAPI/internal/web" // <-- Import the new web package
	"github.com/gorilla/mux"
)

// Env holds application-wide dependencies that are injected into handlers.
type Env struct {
	BookRepo   repository.BookRepository
	UserRepo   repository.UserRepository
	AuthorRepo repository.AuthorRepository
	JWTSecret  string
}

// PaginatedBooksResponse is the structure for paginated book list responses.
type PaginatedBooksResponse struct {
	Metadata map[string]interface{} `json:"metadata"`
	Data     []models.Book          `json:"data"`
}

// GetBooksHandler handles fetching books with pagination, filtering, and sorting.
func (e *Env) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	// ... (Pagination, Filtering, Sorting logic remains the same)
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	var filter repository.BookFilter
	if title := r.URL.Query().Get("title"); title != "" {
		filter.Title = &title
	}
	if author := r.URL.Query().Get("author"); author != "" {
		filter.Author = &author
	}
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	books, totalRecords, err := e.BookRepo.Search(filter, limit, offset, sort, order)
	if err != nil {
		log.Printf("Handler error searching books: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if books == nil {
		books = []models.Book{}
	}

	metadata := map[string]interface{}{
		"current_page":  page,
		"page_size":     limit,
		"total_records": totalRecords,
		"total_pages":   int(math.Ceil(float64(totalRecords) / float64(limit))),
	}

	response := PaginatedBooksResponse{
		Metadata: metadata,
		Data:     books,
	}

	web.RespondWithJSON(w, http.StatusOK, response)
}

// CreateBookHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(newBook); err != nil {
		errors := validationErrors(err)
		web.RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errors})
		return
	}

	_, err := e.AuthorRepo.GetByID(newBook.AuthorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusBadRequest, "Author with the specified ID does not exist")
		} else {
			log.Printf("Handler error checking author existence: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	id, err := e.BookRepo.Create(newBook)
	if err != nil {
		log.Printf("Handler error creating book: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	createdBook, err := e.BookRepo.GetByID(id)
	if err != nil {
		log.Printf("Handler error fetching created book: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, createdBook)
}

// GetBookHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	book, err := e.BookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Book not found")
		} else {
			log.Printf("Handler error getting book by ID: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	web.RespondWithJSON(w, http.StatusOK, book)
}

// UpdateBookHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(updatedBook); err != nil {
		errors := validationErrors(err)
		web.RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errors})
		return
	}

	_, err := e.AuthorRepo.GetByID(updatedBook.AuthorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusBadRequest, "Author with the specified ID does not exist")
		} else {
			log.Printf("Handler error checking author existence: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	err = e.BookRepo.Update(id, updatedBook)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Book not found")
		} else {
			log.Printf("Handler error updating book: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Failed to update book")
		}
		return
	}

	finalBook, err := e.BookRepo.GetByID(id)
	if err != nil {
		log.Printf("Handler error fetching updated book: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	web.RespondWithJSON(w, http.StatusOK, finalBook)
}

// DeleteBookHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := e.BookRepo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Book not found")
		} else {
			log.Printf("Handler error deleting book: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Failed to delete book")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
