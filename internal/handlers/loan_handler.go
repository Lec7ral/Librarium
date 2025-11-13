// Package handlers contains the HTTP handlers for the application.
// This file contains the handlers for loan-related actions.
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

// ... (LoanRequest struct and CreateLoanHandler remain the same)
type LoanRequest struct {
	BookID int64 `json:"book_id" validate:"required"`
}

func (e *Env) CreateLoanHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(web.UserContextKey).(*models.User)
	if !ok {
		web.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve user from context")
		return
	}
	userID := user.ID

	var req LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		errors := validationErrors(err)
		web.RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errors})
		return
	}

	err := e.LoanRepo.CreateLoan(req.BookID, userID)
	if err != nil {
		if err.Error() == "no stock available" {
			web.RespondWithError(w, http.StatusConflict, "No stock available for this book.")
		} else if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Book not found.")
		} else {
			log.Printf("Handler error creating loan: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Failed to process loan.")
		}
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Book loaned successfully."})
}

// ... (ReturnLoanHandler and GetMyLoansHandler remain the same)
func (e *Env) ReturnLoanHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid loan ID")
		return
	}

	err = e.LoanRepo.ReturnLoan(loanID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusNotFound, "Loan not found")
		} else if err.Error() == "book already returned" {
			web.RespondWithError(w, http.StatusConflict, "Book has already been returned")
		} else {
			log.Printf("Handler error returning loan: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Failed to process return")
		}
		return
	}

	web.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book returned successfully."})
}

func (e *Env) GetMyLoansHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(web.UserContextKey).(*models.User)
	if !ok {
		web.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve user from context")
		return
	}

	loans, err := e.LoanRepo.GetActiveLoansByUserID(user.ID)
	if err != nil {
		log.Printf("Handler error getting user loans: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve loans")
		return
	}

	if loans == nil {
		loans = []models.Loan{}
	}

	web.RespondWithJSON(w, http.StatusOK, loans)
}

// @Summary      List all loans (Admin)
// @Description  Get a list of all loans in the system. Can be filtered by status.
// @Tags         Loans
// @Accept       json
// @Produce      json
// @Param        status   query     string  false  "Filter by loan status. Allowed values: active, returned"
// @Success      200      {array}   models.Loan
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /loans [get]
func (e *Env) GetAllLoansHandler(w http.ResponseWriter, r *http.Request) {
	// Create a filter from the query parameters.
	var filter repository.LoanFilter
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}

	// Delegate the query to the repository.
	loans, err := e.LoanRepo.SearchLoans(filter)
	if err != nil {
		log.Printf("Handler error searching loans: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve loans")
		return
	}

	if loans == nil {
		loans = []models.Loan{}
	}

	web.RespondWithJSON(w, http.StatusOK, loans)
}
