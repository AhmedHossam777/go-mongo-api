package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
)

func RegisterUserRoutes(
	router *http.ServeMux, userHandler *handlers.UserHandler,
) {
	const basePath = "/api/v1/users"
	router.HandleFunc("POST "+basePath, userHandler.CreateUser)
	router.HandleFunc("GET "+basePath, userHandler.GetAllUsers)
	router.HandleFunc("GET "+basePath+"/{id}", userHandler.GetOneUser)
	router.HandleFunc("PATCH "+basePath+"/{id}", userHandler.UpdateUser)
	router.HandleFunc("DELETE "+basePath+"/{id}", userHandler.DeleteUser)

	protected := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{"GET", "/api/v1/users/me", userHandler.GetMe},
	}

	for _, route := range protected {
		router.Handle(route.method+" "+route.path,
			middlewares.AuthMiddleware(route.handler))
	}

	// Future user-related endpoints could include:
	// router.HandleFunc("POST /users/login", handler.Login)
	// router.HandleFunc("POST /users/logout", handler.Logout)
	// router.HandleFunc("POST /users/refresh-token", handler.RefreshToken)
	// router.HandleFunc("GET /users/{id}/enrolled-courses", handler.GetUserCourses)
	// router.HandleFunc("PATCH /users/{id}/password", handler.ChangePassword)
	// router.HandleFunc("POST /users/forgot-password", handler.ForgotPassword)
	// router.HandleFunc("POST /users/reset-password", handler.ResetPassword)

	// Keeping authentication and user profile management routes together
	// makes it easier to implement features like role-based access control
	// or audit logging for user-related actions
}
