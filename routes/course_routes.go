package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
)

func RegisterCourseRoutes(
	router *http.ServeMux, courseHandler *handlers.CourseHandler,
) {
	var basePath = "/api/v1/courses"
	router.HandleFunc("GET "+basePath, courseHandler.GetAllCourses)
	router.HandleFunc("GET "+basePath+"/{id}", courseHandler.GetOneCourse)

	protected := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{
			method:  "POST",
			path:    basePath,
			handler: courseHandler.CreateCourse,
		},
		{
			method:  "PATCH",
			path:    basePath + "/{id}",
			handler: courseHandler.UpdateCourse,
		},
		{
			method:  "DELETE",
			path:    basePath + "/{id}",
			handler: courseHandler.DeleteOneCourse,
		},
	}

	for _, route := range protected {
		router.Handle(route.method+" "+route.path,
			middlewares.AuthMiddleware(route.handler))
	}

	//? admin only routes
	router.Handle("DELETE "+basePath+"/drop", middlewares.AuthMiddleware(
		middlewares.RoleMiddleware("admin")(http.HandlerFunc(courseHandler.Drop)),
	))
	// As your application grows, you might add more course-related endpoints here:
	// router.HandleFunc("GET /courses/{id}/students", handler.GetCourseStudents)
	// router.HandleFunc("POST /courses/{id}/enroll", handler.EnrollInCourse)
	// router.HandleFunc("GET /courses/{id}/materials", handler.GetCourseMaterials)
	// router.HandleFunc("POST /courses/{id}/reviews", handler.CreateCourseReview)
}
