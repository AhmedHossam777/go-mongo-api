package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, ErrorResponse{Message: message})
}

func SendSuccess(
	w http.ResponseWriter, status int, message string, data interface{},
) {
	SendJSON(w, status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
