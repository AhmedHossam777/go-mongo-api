package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthService interface {
	Register(ctx context.Context, registerDto dto.RegisterDto, r *http.Request) (
		*dto.AuthResponse, error,
	)
	Login(ctx context.Context, loginDto dto.LoginDto, r *http.Request) (
		*dto.AuthResponse, error,
	)
	RefreshTokens(
		ctx context.Context, refreshToken string, r *http.Request,
	) (*dto.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetActiveSessions(
		ctx context.Context, userID primitive.ObjectID,
	) ([]map[string]interface{}, error)
}

type authService struct {
	userService      UserService
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(
	userService UserService, refreshTokenRepo repository.RefreshTokenRepository,
) AuthService {
	return &authService{
		userService: userService, refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *authService) Register(
	ctx context.Context, registerDto dto.RegisterDto, r *http.Request,
) (*dto.AuthResponse, error) {
	user, _ := s.userService.GetUserByEmail(ctx, registerDto.Email)
	if user != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := helpers.HashPassword(registerDto.Password)
	if err != nil {
		return nil, err
	}

	userModel := &models.User{
		Name:      registerDto.Name,
		Email:     registerDto.Email,
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.userService.CreateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	tokenPairs, err := s.createTokenPair(createdUser, r)
	if err != nil {
		return nil, err
	}

	authResponse := &dto.AuthResponse{
		Token: tokenPairs,
		User: dto.UserResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
		},
	}

	return authResponse, nil
}

func (s *authService) Login(
	ctx context.Context, loginDto dto.LoginDto, r *http.Request,
) (*dto.AuthResponse, error) {

	existedUser, err := s.userService.GetUserByEmail(ctx, loginDto.Email)
	if errors.Is(err, mongo.ErrNoDocuments) || existedUser == nil {
		return nil, ErrInvalidCredentials
	}

	isCorrect := helpers.CheckPassword(existedUser.Password, loginDto.Password)
	if !isCorrect {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := s.createTokenPair(existedUser, r)

	if err != nil {
		return nil, err
	}

	authResponse := &dto.AuthResponse{
		Token: tokenPair,
		User: dto.UserResponse{
			ID:    existedUser.ID,
			Name:  existedUser.Name,
			Email: existedUser.Email,
		},
	}

	return authResponse, nil
}

func (s *authService) RefreshTokens(
	ctx context.Context, refreshToken string, r *http.Request,
) (*dto.TokenPair, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	matchedToken, err := s.findValidRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenRepo.RevokeToken(ctx, matchedToken.ID)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.GetOneUser(ctx, matchedToken.UserId.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.createTokenPair(user, r)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	matchedToken, err := s.findValidRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	return s.refreshTokenRepo.RevokeToken(ctx, matchedToken.ID)
}
func (s *authService) GetActiveSessions(
	ctx context.Context, userID primitive.ObjectID,
) ([]map[string]interface{}, error) {
	tokens, err := s.refreshTokenRepo.FindActiveTokensByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("error fetching sessions")
	}

	sessions := make([]map[string]interface{}, 0, len(tokens))
	for _, token := range tokens {
		sessions = append(
			sessions, map[string]interface{}{
				"id":        token.ID.Hex(),
				"userAgent": token.UserAgent,
				"ipAddress": token.IPAddress,
				"createdAt": token.CreatedAt,
				"expiresAt": token.ExpiresAt,
			},
		)
	}

	return sessions, nil
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the list
		return strings.Split(forwarded, ",")[0]
	}
	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (s *authService) createTokenPair(user *models.User, r *http.Request) (
	*dto.TokenPair, error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accessToken, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	plainRefreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := helpers.HashRefreshToken(plainRefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenDoc := models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserId:    user.ID,
		Token:     refreshToken,
		ExpiresAt: helpers.GetRefreshTokenExpiry(),
		CreatedAt: time.Now(),
		Revoked:   false,
		RevokedAt: nil,
		UserAgent: r.UserAgent(),
		IPAddress: getClientIP(r),
	}

	err = s.refreshTokenRepo.Create(ctx, &refreshTokenDoc)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: plainRefreshToken,
	}, nil
}

func (s *authService) findValidRefreshToken(
	ctx context.Context, plainToken string,
) (*models.RefreshToken, error) {
	activeTokens, err := s.refreshTokenRepo.FindActiveTokens(ctx)
	if err != nil {
		return nil, err
	}

	for _, token := range activeTokens {
		if helpers.ValidateRefreshToken(plainToken, token.Token) {
			return token, nil
		}
	}

	return nil, errors.New("token not found")
}
