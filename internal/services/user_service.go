package services

import (
	"context"
	"errors"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrUsernameRequired     = errors.New("username is required")
	ErrPasswordRequired     = errors.New("password is required")
	ErrEmailRequired        = errors.New("email is required")
	ErrNoFieldsToUpdateUser = errors.New("no fields to update")
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetOneUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *dto.UpdateUserDto) (
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
	ctx context.Context, id string, updateUserDto *dto.UpdateUserDto,
) (*models.User, error) {

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	var update = bson.M{}
	if updateUserDto.Name != nil {
		update["Name"] = updateUserDto.Name
	}
	if updateUserDto.Email != nil {
		update["Email"] = updateUserDto.Email
	}
	if updateUserDto.Password != nil {
		update["Password"] = updateUserDto.Password
	}

	updatedUser, err := s.repo.UpdateOneUser(ctx, objId, bson.M{"$set": update})

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
