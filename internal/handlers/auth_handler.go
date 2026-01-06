package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(
	authService services.AuthService,
) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterDto true "User registration details"
// @Success 201 {object} dto.AuthResponse "User registered successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or user already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var registerDto dto.RegisterDto

	err := json.NewDecoder(r.Body).Decode(&registerDto)
	defer r.Body.Close()

	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error(),
		)
		return
	}

	validationErr := helpers.ValidateStruct(registerDto)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	authResponse, err := h.authService.Register(ctx, registerDto, r)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			RespondWithError(
				w, http.StatusBadRequest,
				"user already exist",
			)
			return
		}
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while register, "+err.Error(),
		)
		return
	}

	RespondWithJSON(w, http.StatusCreated, authResponse)
}

// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginDto true "User login credentials"
// @Success 200 {object} dto.AuthResponse "User logged in successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginDto dto.LoginDto
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error(),
		)
		return
	}

	validationErr := helpers.ValidateStruct(loginDto)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	authResponse, err := h.authService.Login(ctx, loginDto, r)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while login, "+err.Error(),
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, authResponse)
}

// @Summary Refresh access token
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenInput true "Refresh token"
// @Success 200 {object} dto.TokenPair "New token pair generated"
// @Failure 401 {object} map[string]string "Unauthorized - invalid refresh token"
// @Router /auth/refresh-tokens [post]
func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var refreshTokenInput *dto.RefreshTokenInput
	err := json.NewDecoder(r.Body).Decode(&refreshTokenInput)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error(),
		)
		return
	}

	validationErr := helpers.ValidateStruct(refreshTokenInput)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	tokenParis, err := h.authService.RefreshTokens(
		ctx, refreshTokenInput.RefreshToken, r,
	)

	if err != nil {
		RespondWithError(
			w, http.StatusUnauthorized, "error while refresh tokens, "+err.Error(),
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, tokenParis)
}

// @Summary Logout user
// @Description Logout user and invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenInput true "Refresh token to invalidate"
// @Success 200 {object} map[string]string "Logged out successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var refreshTokenInput *dto.RefreshTokenInput
	err := json.NewDecoder(r.Body).Decode(&refreshTokenInput)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while decoding request body: "+err.Error(),
		)
		return
	}

	validationErr := helpers.ValidateStruct(refreshTokenInput)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	err = h.authService.Logout(ctx, refreshTokenInput.RefreshToken)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Error while logging out: "+err.Error(),
		)
		return
	}

	RespondWithJSON(
		w, http.StatusOK, map[string]string{"message": "Logged out successfully"},
	)
}

// @Summary Get active sessions
// @Description Get all active sessions for the authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of active sessions"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/active-sessions [get]
func (h *AuthHandler) GetActiveSessions(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(
			w, http.StatusUnauthorized,
			"user id not found",
		)
		return
	}

	mongoUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError, "invalid user id in the context",
		)
	}

	sessions, err := h.authService.GetActiveSessions(ctx, mongoUserId)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError, "error getting active session",
		)
		return
	}

	RespondWithJSON(
		w, http.StatusOK, map[string]interface{}{
			"activeSessions": sessions,
			"count":          len(sessions),
		},
	)
}
