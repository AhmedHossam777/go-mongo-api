package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
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

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default to 10
	}

	users, totalCount, err := h.service.GetAllUsers(ctx, int64(page),
		int64(pageSize))

	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError,
			"error while fetching all user, "+err.Error())
		return
	}

	hasMore := false
	count := pageSize
	if int(totalCount) > (page * pageSize) {
		hasMore = true
	} else {
		count = pageSize - (page*pageSize - int(totalCount))
	}

	PaginationResponse(w, http.StatusOK, users, page, count,
		totalCount,
		hasMore)
}

func (h *UserHandler) GetMe(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.Context().Value("userId").(string)
	fmt.Println("userId: ", userId)
	user, err := h.service.GetOneUser(ctx, userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while getting one user")
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
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

	if createUserDto.Role == "" {
		createUserDto.Role = "user"
	}

	user := &models.User{
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  hashedPassword,
		Role:      createUserDto.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := h.service.CreateUser(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			RespondWithError(w, http.StatusBadRequest,
				"user already exist")
			return
		}
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
			"error while deleting one user, "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}

func (h *UserHandler) DropUserCollection(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.service.DropUserCollection(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while deleting all users, "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
