package middleware

import (
	"github.com/Lec7ral/fullAPI/internal/models"
	"github.com/Lec7ral/fullAPI/internal/web"
	"net/http"
)

func RoleRequiredMiddleware(rquieredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(web.UserContextKey).(*models.User)
			if !ok {
				web.RespondWithError(w, http.StatusUnauthorized, "User Not found in context")
				return
			}
			if user.Role != rquieredRole {
				web.RespondWithError(w, http.StatusForbidden, "You don't have permission to access this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
