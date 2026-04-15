package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
)

func RegisterBlogRoutes(router *http.ServeMux, blogHandler *handlers.BlogHandler) {
	var basePath = "/api/v1/blogs"

	router.HandleFunc("GET "+basePath, blogHandler.GetAllBlogs)
	router.HandleFunc("GET "+basePath+"/{id}", blogHandler.GetOneBlog)
	router.HandleFunc("GET "+basePath+"/{id}/comments", blogHandler.GetComments)

	protected := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{
			method:  "POST",
			path:    basePath,
			handler: blogHandler.CreateBlog,
		},
		{
			method:  "PUT",
			path:    basePath + "/{id}",
			handler: blogHandler.UpdateBlog,
		},
		{
			method:  "DELETE",
			path:    basePath + "/{id}",
			handler: blogHandler.DeleteOneBlog,
		},
		{
			method:  "POST",
			path:    basePath + "/{id}/comments",
			handler: blogHandler.AddComment,
		},
		{
			method:  "GET",
			path:    basePath + "/my",
			handler: blogHandler.GetBlogsByAuthor,
		},
	}

	for _, route := range protected {
		router.Handle(route.method+" "+route.path,
			middlewares.AuthMiddleware(route.handler))
	}
}
