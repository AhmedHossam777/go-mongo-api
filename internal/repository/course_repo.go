package repository

import (
	"context"
	"errors"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) (*models.Course, error)
	FindAll(ctx context.Context) ([]models.Course, error)
	FindOne(ctx context.Context, courseId primitive.ObjectID) (
		*models.Course, error,
	)
	UpdateOne(
		ctx context.Context, courseId primitive.ObjectID, update bson.M,
	) (
		*models.Course, error,
	)
	DeleteOne(ctx context.Context, courseId primitive.ObjectID) error
}

type courseRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewCourseRepo(db *mongo.Database) CourseRepository {
	return &courseRepository{
		collection: db.Collection("courses"),
		timeout:    10 * time.Second}
}

func (r *courseRepository) Create(
	ctx context.Context, course *models.Course,
) (*models.Course, error) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	course.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *courseRepository) FindAll(ctx context.Context) (
	[]models.Course, error,
) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var courses []models.Course
	err = cursor.All(ctx, &courses)
	if err != nil {
		return nil, err
	}

	// if courses table is empty return empty slice not nil
	if courses == nil {
		courses = []models.Course{}
	}

	return courses, nil
}

func (r *courseRepository) FindOne(
	ctx context.Context, id primitive.ObjectID,
) (*models.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var course *models.Course

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r courseRepository) UpdateOne(
	ctx context.Context, id primitive.ObjectID, update bson.M,
) (*models.Course, error) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedCourse models.Course

	err := r.collection.FindOneAndUpdate(ctx, filter, update,
		opts).Decode(&updatedCourse)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	return &updatedCourse, nil
}

func (r *courseRepository) DeleteOne(
	ctx context.Context, id primitive.ObjectID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
