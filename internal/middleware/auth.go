// Package middleware provides HTTP middleware functions for the application.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Lec7ral/fullAPI/internal/repository"
	"github.com/Lec7ral/fullAPI/internal/web"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is a constructor that takes dependencies (UserRepo, JWTSecret)
// and returns a middleware handler. This is the standard way to inject dependencies
// into middleware without causing circular imports.
func AuthMiddleware(userRepo repository.UserRepository, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				web.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				web.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Use the username from the token to fetch the full user object.
			username := claims.Subject
			user, err := userRepo.GetByUsername(username)
			if err != nil {
				web.RespondWithError(w, http.StatusUnauthorized, "User not found")
				return
			}

			// Store the user object in the context using the exported key.
			ctx := context.WithValue(r.Context(), web.UserContextKey, user)

			// Call the next handler with the enriched context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
