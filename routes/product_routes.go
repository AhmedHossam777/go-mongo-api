package routes

import (
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/middlewares"
)

func RegisterProductRoutes(
	router *http.ServeMux, productHandler *handlers.ProductHandler,
) {
	var basePath = "/api/v1/products"
	router.HandleFunc("GET "+basePath, productHandler.GetAllProducts)
	router.HandleFunc("GET "+basePath+"/{id}", productHandler.GetOneProduct)
	router.HandleFunc("POST "+basePath, productHandler.CreateProduct)
	router.HandleFunc("PATCH "+basePath+"/{id}", productHandler.UpdateProduct)
	router.HandleFunc("DELETE "+basePath+"/{id}", productHandler.DeleteOneProduct)

	// protected := []struct {
	// 	method  string
	// 	path    string
	// 	handler http.HandlerFunc
	// }{
	// 	{
	// 		method:  "POST",
	// 		path:    basePath,
	// 		handler: productHandler.CreateProduct,
	// 	},
	// 	{
	// 		method:  "PATCH",
	// 		path:    basePath + "/{id}",
	// 		handler: productHandler.UpdateProduct,
	// 	},
	// 	{
	// 		method:  "DELETE",
	// 		path:    basePath + "/{id}",
	// 		handler: productHandler.DeleteOneProduct,
	// 	},
	// }
	//
	// for _, route := range protected {
	// 	router.Handle(
	// 		route.method+" "+route.path,
	// 		middlewares.AuthMiddleware(route.handler),
	// 	)
	// }

	// ? admin only routes
	router.Handle(
		"DELETE "+basePath+"/drop", middlewares.AuthMiddleware(
			middlewares.RoleMiddleware("admin")(http.HandlerFunc(productHandler.Drop)),
		),
	)
}
