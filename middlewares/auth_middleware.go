package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.RespondWithError(w, http.StatusUnauthorized,
				"Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			handlers.RespondWithError(w, http.StatusUnauthorized,
				"Authorization header format must be 'Bearer <token>'")
			return
		}
		tokenString := parts[1]

		claim, err := helpers.ValidateToken(tokenString)
		if err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized,
				"Invalid or expired token:  "+err.Error())
			return
		}

		//Add claims to request context
		// This allows handlers to access user info
		ctx := context.WithValue(r.Context(), "userId", claim.UserId)
		ctx = context.WithValue(ctx, "userEmail", claim.Email)
		ctx = context.WithValue(ctx, "userRole", claim.Role)

		//Call the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
