package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
)

func RegisterAuthRouts(
	router *http.ServeMux, authHandler *handlers.AuthHandler,
) {

	const basePath = "/api/v1/auth"
	router.HandleFunc("POST "+basePath+"/register", authHandler.Register)
	router.HandleFunc("POST "+basePath+"/login", authHandler.Login)
	router.HandleFunc(
		"POST "+basePath+"/refresh-tokens", authHandler.RefreshTokens,
	)
	router.HandleFunc("POST "+basePath+"/logout", authHandler.Logout)

	router.Handle(
		"GET "+basePath+"/active-sessions",
		middlewares.AuthMiddleware(http.HandlerFunc(authHandler.GetActiveSessions)),
	)
}
