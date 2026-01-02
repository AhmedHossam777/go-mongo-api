package middlewares

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
)

func RoleMiddleware(allowedStr ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user role
			userRole, ok := r.Context().Value("userRole").(string)
			if !ok {
				handlers.RespondWithError(w, http.StatusUnauthorized,
					"user role not found")
				return
			}

			roleAllowed := false
			for _, role := range allowedStr {
				if role == userRole {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				handlers.RespondWithError(w, http.StatusForbidden,
					"You don't have permission to access this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
