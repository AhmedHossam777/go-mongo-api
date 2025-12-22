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

	//userRepo := repository.NewUserRepo(db)

	port := cnfg.Port
	if port == "" {
		port = "3000"
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /", serverHome)
	router.HandleFunc("POST /courses", courseHandler.CreateCourse)
	router.HandleFunc("GET /courses", courseHandler.GetAllCourses)
	router.HandleFunc("GET /courses/{id}", courseHandler.GetOneCourse)
	router.HandleFunc("PATCH /courses/{id}", courseHandler.UpdateCourse)
	router.HandleFunc("DELETE /courses/{id}", courseHandler.DeleteOneCourse)

	// Start server
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📚 API available at http://localhost:%s/courses\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func serverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, in my first golang server"})
}
