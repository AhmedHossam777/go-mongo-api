package services

import (
	"context"
	"errors"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
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
