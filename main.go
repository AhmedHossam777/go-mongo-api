package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AhmedHossam777/go-mongo/config"
	"github.com/AhmedHossam777/go-mongo/controllers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found!!")
	}

	config.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /", serverHome)
	router.HandleFunc("POST /courses", controllers.CreateCourse)
	// Start server
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📚 API available at http://localhost:%s/courses\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func serverHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, in my first golang server"})
}
