// Package middleware provides HTTP middleware functions for the application.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Lec7ral/fullAPI/internal/web" // <-- Import the new web package
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is a constructor that takes a JWT secret key and returns a middleware handler.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Use the helper from the web package.
				web.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				// Use the helper from the web package.
				web.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), "username", claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
