package services

import (
	"context"
	"errors"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrUsernameRequired     = errors.New("user name is required")
	ErrPasswordRequired     = errors.New("user name is required")
	ErrEmailRequired        = errors.New("user name is required")
	ErrNoFieldsToUpdateUser = errors.New("no fields to update")
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetOneUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (
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

func (s *userService) CreateUser(
	ctx context.Context, user *models.User,
) (*models.User, error) {
	if user.Name == "" {
		return nil, ErrUsernameRequired
	}
	if user.Email == "" {
		return nil, ErrEmailRequired
	}
	if user.Password == "" {
		return nil, ErrPasswordRequired
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *userService) GetAllUsers(
	ctx context.Context,
) ([]models.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *userService) GetOneUser(
	ctx context.Context, id string,
) (*models.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	user, err := s.repo.GetOneUser(ctx, objId)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(
	ctx context.Context, id string, user *models.User,
) (*models.User, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	if user.Role == "" && user.Name == "" && user.Email == "" && user.Password == "" {
		return nil, ErrNoFieldsToUpdateUser
	}

	updatedUser, err := s.repo.UpdateOneUser(ctx, objId, user)

	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteUser(
	ctx context.Context, id string,
) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidUserID
	}
	err = s.repo.DeleteOneUser(ctx, objId)

	if err == mongo.ErrNoDocuments {
		return ErrCourseNotFound
	}
	return nil
}
