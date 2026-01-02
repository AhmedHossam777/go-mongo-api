package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
)

func RegisterAuthRouts(
	router *http.ServeMux, authHandler *handlers.AuthHandler,
) {

	const basePath = "/api/v1/auth"
	router.HandleFunc("POST "+basePath+"/register", authHandler.Register)
	router.HandleFunc("POST "+basePath+"/login", authHandler.Login)
}
