package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/AhmedHossam777/go-mongo/config"
	"go.mongodb.org/mongo-driver/mongo"
)

//! some helper functions

func getCourseCollection() *mongo.Collection {
	return config.DB.Collection("courses")
}

func sendJson(
	w http.ResponseWriter, status int, data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		sendErr(w, 500, "error while encoding json response")
	}
}

func sendErr(
	w http.ResponseWriter, status int, message string,
) {
	sendJson(w, status, map[string]string{"message": message})
}
