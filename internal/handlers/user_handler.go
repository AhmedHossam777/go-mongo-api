package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

func (h *UserHandler) GetAllUsers(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := h.service.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		SendError(w, http.StatusInternalServerError,
			"error while fetching all user"+err.Error())
		return
	}

	SendSuccess(w, http.StatusOK, "User fetched successfully", users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Decode into SignUpInput, not User
	var input models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&input)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		SendError(w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error())
		return
	}

	// Validate input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		SendError(w, http.StatusBadRequest,
			"Name, email, and password are required")
		return
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		SendError(w, http.StatusInternalServerError,
			"Error while hashing the password")
		return
	}

	// Create User from input
	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := h.service.CreateUser(ctx, user)
	if err != nil {
		SendError(w, http.StatusInternalServerError,
			"Error while creating new user: "+err.Error())
		return
	}

	// Return response without password
	SendSuccess(w, http.StatusCreated, "User created successfully",
		createdUser.ToResponse())
}

func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")
	user, err := h.service.GetOneUser(ctx, userId)
	if err != nil {
		SendError(w, http.StatusInternalServerError, "error while getting one user")
		return
	}

	SendSuccess(w, http.StatusOK, "User fetched successfully", user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userData *models.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		log.Println(err)
		SendError(w, http.StatusInternalServerError,
			"Error while decoding request body")
		return
	}

	userId := r.PathValue("id")
	updatedUser, err := h.service.UpdateUser(ctx, userId, userData)
	if err != nil {
		SendError(w, http.StatusInternalServerError,
			"error while updating one user")
		return
	}

	SendSuccess(w, http.StatusOK, "User updated successfully", updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")

	err := h.service.DeleteUser(ctx, userId)
	if err != nil {
		SendError(w, http.StatusInternalServerError,
			"error while deleting one user")
	}

	SendSuccess(w, http.StatusOK, "User Deleted successfully", nil)
}
