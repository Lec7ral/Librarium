// Package handlers contains the HTTP handlers for the application.
// This file focuses on user registration and login.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Lec7ral/fullAPI/internal/models"
	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/Lec7ral/fullAPI/internal/web" // <-- Import the new web package
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Credentials defines the structure for user login and registration requests.
type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterUserHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(creds); err != nil {
		errors := validationErrors(err)
		web.RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": errors})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user := models.User{Username: creds.Username}
	err = e.UserRepo.Create(user, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUsernameExists) {
			web.RespondWithError(w, http.StatusConflict, "Username already exists")
		} else {
			log.Printf("Handler error creating user: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, nil)
}

// LoginUserHandler now uses the standardized JSON response helpers from the web package.
func (e *Env) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := e.UserRepo.GetByUsername(creds.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			web.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		} else {
			log.Printf("Handler error getting user: %v", err)
			web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		web.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(e.JWTSecret))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		web.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	web.RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
