package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/config"
	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"github.com/AhmedHossam777/go-mongo/internal/services"
	"github.com/AhmedHossam777/go-mongo/routes"

	_ "github.com/AhmedHossam777/go-mongo/docs"
)

// @title Go-MongoDB Course API
// @version 1.0
// @description A RESTful API for managing courses, users, and authentication built with Go and MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cnfg := config.LoadConfig()

	db, err := config.ConnectDB(cnfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	courseRepo := repository.NewCourseRepo(db)
	courseService := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	userRepo := repository.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	refreshTokenRepo := repository.NewRefreshTokenRepo(db)
	authService := services.NewAuthService(userService, refreshTokenRepo)
	authHandler := handlers.NewAuthHandler(authService)

	productRepo := repository.NewProductRepo(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	blogRepo := repository.NewBlogRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	blogService := services.NewBlogService(blogRepo, commentRepo)
	s3Service, err := services.NewS3Service()
	if err != nil {
		log.Fatal("Failed to initialize S3 service:", err)
	}
	blogHandler := handlers.NewBlogHandler(blogService, s3Service)

	port := cnfg.Port
	if port == "" {
		port = "8080"
	}

	router := routes.SetupRoutes(
		userHandler, courseHandler, authHandler, productHandler, blogHandler,
	)

	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║       Go-MongoDB Course API Server                ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📚 API v1 Base URL: http://localhost:%s/api/v1\n", port)
	fmt.Printf("📖 Courses API: http://localhost:%s/api/v1/courses\n", port)
	fmt.Printf("📦 Products API: http://localhost:%s/api/v1/products\n", port)
	fmt.Printf("📝 Blogs API: http://localhost:%s/api/v1/blogs\n", port)
	fmt.Printf("👥 Users API: http://localhost:%s/api/v1/users\n", port)
	fmt.Printf("🔐 Auth API: http://localhost:%s/api/v1/auth\n", port)
	fmt.Printf("📝 Swagger Docs: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("💚 Health Check: http://localhost:%s/health\n", port)
	fmt.Println("════════════════════════════════════════════════════")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
