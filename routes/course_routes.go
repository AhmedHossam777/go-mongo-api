package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
)

func RegisterCourseRoutes(
	router *http.ServeMux, courseHandler *handlers.CourseHandler,
) {
	var basePath = "/api/v1/courses"
	router.HandleFunc("POST "+basePath, courseHandler.CreateCourse)
	router.HandleFunc("GET "+basePath, courseHandler.GetAllCourses)
	router.HandleFunc("GET "+basePath+"/{id}", courseHandler.GetOneCourse)
	router.HandleFunc("PATCH "+basePath+"/{id}", courseHandler.UpdateCourse)
	router.HandleFunc("DELETE "+basePath+"/{id}", courseHandler.DeleteOneCourse)

	// As your application grows, you might add more course-related endpoints here:
	// router.HandleFunc("GET /courses/{id}/students", handler.GetCourseStudents)
	// router.HandleFunc("POST /courses/{id}/enroll", handler.EnrollInCourse)
	// router.HandleFunc("GET /courses/{id}/materials", handler.GetCourseMaterials)
	// router.HandleFunc("POST /courses/{id}/reviews", handler.CreateCourseReview)
}
