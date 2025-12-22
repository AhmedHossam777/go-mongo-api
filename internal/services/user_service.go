package services

import (
	"context"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetOneUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user models.User) (
		*models.User, error,
	)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
