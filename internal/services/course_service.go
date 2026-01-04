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
	ErrCourseNotFound  = errors.New("course not found")
	ErrInvalidCourseID = errors.New("invalid course ID")
)

type CourseService interface {
	CreateCourse(ctx context.Context, courseDto *dto.CreateCourseDto) (
		*models.Course, error,
	)
	GetAllCourses(ctx context.Context, page, pageSize int64) (
		[]models.Course, int, error,
	)
	GetCourseByID(ctx context.Context, id string) (*models.Course, error)
	UpdateCourse(
		ctx context.Context, id string, updateCourseDto *dto.UpdateCourseDto,
	) (*models.Course, error)
	DeleteCourse(ctx context.Context, id string) error
	Drop(ctx context.Context) error
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) CreateCourse(
	ctx context.Context, courseDto *dto.CreateCourseDto,
) (*models.Course, error) {

	course := &models.Course{
		CourseName: courseDto.CourseName,
		Price:      courseDto.Price,
		Author:     courseDto.AuthorID,
	}
	return s.repo.Create(ctx, course)
}

func (s *courseService) GetAllCourses(
	ctx context.Context, page, pageSize int64,
) (
	[]models.Course, int, error,
) {
	return s.repo.FindAll(ctx, page, pageSize)
}

func (s *courseService) GetCourseByID(ctx context.Context, id string) (
	*models.Course, error,
) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidCourseID
	}

	course, err := s.repo.FindOne(ctx, objectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s *courseService) UpdateCourse(
	ctx context.Context, id string, updateCourseDto *dto.UpdateCourseDto,
) (*models.Course, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidCourseID
	}

	var update = bson.M{}
	if updateCourseDto.CourseName != nil {
		update["CourseName"] = *updateCourseDto.CourseName
	}
	if updateCourseDto.Price != nil {
		update["Price"] = *updateCourseDto.Price
	}

	updateCourse, err := s.repo.UpdateOne(ctx, objectId, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	return updateCourse, nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidCourseID
	}

	err = s.repo.DeleteOne(ctx, objectId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrCourseNotFound
	}
	return nil
}

func (s *courseService) Drop(ctx context.Context) error {
	err := s.repo.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}
