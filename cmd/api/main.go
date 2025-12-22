package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AhmedHossam777/go-mongo/internal/config"
	"github.com/AhmedHossam777/go-mongo/internal/handlers"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"github.com/AhmedHossam777/go-mongo/internal/services"
	"github.com/AhmedHossam777/go-mongo/routes"
)

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

	port := cnfg.Port
	if port == "" {
		port = "3000"
	}

	router := routes.SetupRoutes(userHandler, courseHandler)

	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║       Go-MongoDB Course API Server                ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📚 API v1 Base URL: http://localhost:%s/api/v1\n", port)
	fmt.Printf("📖 Courses API: http://localhost:%s/api/v1/courses\n", port)
	fmt.Printf("👥 Users API: http://localhost:%s/api/v1/users\n", port)
	fmt.Printf("💚 Health Check: http://localhost:%s/health\n", port)
	fmt.Println("════════════════════════════════════════════════════")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func serverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, in my first golang server"})
}
