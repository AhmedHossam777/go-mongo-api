package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/services"
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
