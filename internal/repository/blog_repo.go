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

type BlogRepository interface {
	Create(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	FindAll(ctx context.Context, page int64, pageSize int64) ([]models.Blog, int, error)
	FindOne(ctx context.Context, blogId primitive.ObjectID) (*models.Blog, error)
	UpdateOne(ctx context.Context, blogId primitive.ObjectID, update bson.M) (*models.Blog, error)
	DeleteOne(ctx context.Context, blogId primitive.ObjectID) error
}

type blogRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewBlogRepo(db *mongo.Database) BlogRepository {
	return &blogRepository{
		collection: db.Collection("blogs"),
		timeout:    10 * time.Second,
	}
}

func (r *blogRepository) Create(ctx context.Context, blog *models.Blog) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	blog.ID = primitive.NewObjectID()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (r *blogRepository) FindAll(ctx context.Context, page int64, pageSize int64) ([]models.Blog, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	skip := (page - 1) * pageSize

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(pageSize)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog
	err = cursor.All(ctx, &blogs)
	if err != nil {
		return nil, 0, err
	}

	if blogs == nil {
		blogs = []models.Blog{}
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return blogs, int(totalCount), nil
}

func (r *blogRepository) FindOne(ctx context.Context, id primitive.ObjectID) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var blog models.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (r *blogRepository) UpdateOne(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedBlog models.Blog

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedBlog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}

	return &updatedBlog, nil
}

func (r *blogRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
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
