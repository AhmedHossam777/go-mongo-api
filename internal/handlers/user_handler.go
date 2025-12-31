package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
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
		RespondWithError(w, http.StatusInternalServerError,
			"error while fetching all user"+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var createUserDto dto.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&createUserDto)
	defer r.Body.Close()
	if err != nil {
		RespondWithError(w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error())
		return
	}

	hashedPassword, err := helpers.HashPassword(createUserDto.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"Error while hashing the password")
		return
	}

	validationErr := helpers.ValidateStruct(createUserDto)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	user := &models.User{
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := h.service.CreateUser(ctx, user)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"Error while creating new user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated,
		createdUser.ToResponse())
}

func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")
	user, err := h.service.GetOneUser(ctx, userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while getting one user")
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var updateUserDto dto.UpdateUserDto
	err := json.NewDecoder(r.Body).Decode(&updateUserDto)
	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError,
			"Error while decoding request body")
		return
	}

	validationErr := helpers.ValidateStruct(updateUserDto)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	userId := r.PathValue("id")
	updatedUser, err := h.service.UpdateUser(ctx, userId, &updateUserDto)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while updating one user")
		return
	}

	RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")

	err := h.service.DeleteUser(ctx, userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while deleting one user")
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
