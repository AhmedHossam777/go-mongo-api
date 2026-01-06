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

// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags users
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Paginated list of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
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

	users, totalCount, err := h.service.GetAllUsers(
		ctx, int64(page),
		int64(pageSize),
	)

	if err != nil {
		log.Println(err)
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while fetching all user, "+err.Error(),
		)
		return
	}

	hasMore := false
	count := pageSize
	if int(totalCount) > (page * pageSize) {
		hasMore = true
	} else {
		count = pageSize - (page*pageSize - int(totalCount))
	}

	PaginationResponse(
		w, http.StatusOK, users, page, count,
		totalCount,
		hasMore,
	)
}

// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.User "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/me [get]
func (h *UserHandler) GetMe(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.Context().Value("userId").(string)
	fmt.Println("userId: ", userId)
	user, err := h.service.GetOneUser(ctx, userId)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while getting one user",
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

// @Summary Create a new user
// @Description Create a new user with name, email, password and optional role
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserDto true "User details"
// @Success 201 {object} github_com_AhmedHossam777_go-mongo_internal_models.UserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or user already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var createUserDto dto.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&createUserDto)
	defer r.Body.Close()
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error(),
		)
		return
	}

	hashedPassword, err := helpers.HashPassword(createUserDto.Password)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"Error while hashing the password",
		)
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
			RespondWithError(
				w, http.StatusBadRequest,
				"user already exist",
			)
			return
		}
		RespondWithError(
			w, http.StatusInternalServerError,
			"Error while creating new user: "+err.Error(),
		)
		return
	}

	RespondWithJSON(
		w, http.StatusCreated,
		createdUser.ToResponse(),
	)
}

// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.User "User details"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetOneUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")
	user, err := h.service.GetOneUser(ctx, userId)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while getting one user",
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

// @Summary Update user
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserDto true "Updated user details"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.User "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var updateUserDto dto.UpdateUserDto
	err := json.NewDecoder(r.Body).Decode(&updateUserDto)
	if err != nil {
		log.Println(err)
		RespondWithError(
			w, http.StatusInternalServerError,
			"Error while decoding request body",
		)
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
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while updating one user",
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, updatedUser)
}

// @Summary Delete user
// @Description Delete a user by ID (Admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := r.PathValue("id")

	err := h.service.DeleteUser(ctx, userId)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while deleting one user, "+err.Error(),
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}

// @Summary Drop user collection
// @Description Delete all users from the database (Admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 "All users deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/drop [delete]
func (h *UserHandler) DropUserCollection(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.service.DropUserCollection(ctx)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"error while deleting all users, "+err.Error(),
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
