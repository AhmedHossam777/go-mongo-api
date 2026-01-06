package routes

import (
	"encoding/json"
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(
	userHandler *handlers.UserHandler, courseHandler *handlers.CourseHandler,
	authHandler *handlers.AuthHandler,
) http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("GET /", serverHome)
	router.HandleFunc("GET /health", healthCheck)

	// Swagger documentation
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	RegisterCourseRoutes(router, courseHandler)
	RegisterUserRoutes(router, userHandler)
	RegisterAuthRouts(router, authHandler)

	// Wrap router with CORS middleware
	return middlewares.CORSMiddleware(router)
}

// @Summary Server Home
// @Description Get welcome message and API information
// @Tags general
// @Produce json
// @Success 200 {object} map[string]interface{} "Welcome message with API information"
// @Router / [get]
func serverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message": "Welcome to the Go-MongoDB Course API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": map[string]string{
				"courses": "/courses",
				"users":   "/users",
				"health":  "/health",
			},
		},
	)
}

// @Summary Health Check
// @Description Check if the API is running
// @Tags general
// @Produce json
// @Success 200 {object} map[string]string "API health status"
// @Router /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "healthy",
		},
	)
}
