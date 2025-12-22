package routes

import (
	"encoding/json"
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
)

func SetupRoutes(
	userHandler *handlers.UserHandler, courseHandler *handlers.CourseHandler,
) *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("GET /", serverHome)
	router.HandleFunc("GET /health", healthCheck)

	RegisterCourseRoutes(router, courseHandler)
	RegisterUserRoutes(router, userHandler)

	return router
}

func serverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome to the Go-MongoDB Course API",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]string{
			"courses": "/courses",
			"users":   "/users",
			"health":  "/health",
		},
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
