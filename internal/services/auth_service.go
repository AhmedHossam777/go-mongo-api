package services

import (
	"context"
	"errors"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type AuthService interface {
	Register(ctx context.Context, registerDto dto.RegisterDto) (
		*dto.AuthResponse, error,
	)
	Login(ctx context.Context, loginDto dto.LoginDto) (
		*dto.AuthResponse, error,
	)
}
type authService struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &authService{userService: userService}
}

func (s *authService) Register(
	ctx context.Context, registerDto dto.RegisterDto,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.userService.CreateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	token, err := helpers.GenerateToken(createdUser.ID, createdUser.Email,
		createdUser.Role)

	if err != nil {
		return nil, err
	}

	authResponse := &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
		},
	}

	return authResponse, nil
}

func (s *authService) Login(
	ctx context.Context, loginDto dto.LoginDto,
) (*dto.AuthResponse, error) {

	existedUser, err := s.userService.GetUserByEmail(ctx, loginDto.Email)
	if errors.Is(err, mongo.ErrNoDocuments) || existedUser == nil {
		return nil, ErrInvalidCredentials
	}

	token, err := helpers.GenerateToken(existedUser.ID, existedUser.Email,
		existedUser.Role)

	if err != nil {
		return nil, err
	}

	authResponse := &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    existedUser.ID,
			Name:  existedUser.Name,
			Email: existedUser.Email,
		},
	}

	return authResponse, nil
}
