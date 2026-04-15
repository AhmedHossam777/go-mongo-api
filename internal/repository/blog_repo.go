package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *models.Blog) (*models.Blog, error)
	FindAll(ctx context.Context, page int64, pageSize int64) ([]models.BlogWithAuthor, int, error)
	FindOne(ctx context.Context, blogId primitive.ObjectID) (*models.BlogWithAuthor, error)
	UpdateOne(ctx context.Context, blogId primitive.ObjectID, update bson.M) (*models.BlogWithAuthor, error)
	DeleteOne(ctx context.Context, blogId primitive.ObjectID) error
	GetBlogsByAuthor(ctx context.Context, authorId primitive.ObjectID, page int64, pageSize int64) ([]models.BlogWithAuthor, int, error)
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

var blogAuthorLookup = mongo.Pipeline{
	{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "author_id",
		"foreignField": "_id",
		"as":           "authorArr",
	}}},
	{{Key: "$unwind", Value: bson.M{"path": "$authorArr", "preserveNullAndEmptyArrays": true}}},
	{{Key: "$addFields", Value: bson.M{
		"author": bson.M{
			"id":        "$authorArr._id",
			"name":      "$authorArr.name",
			"email":     "$authorArr.email",
			"role":      "$authorArr.role",
			"createdAt": "$authorArr.created_at",
		},
	}}},
	{{Key: "$project", Value: bson.M{"authorArr": 0, "author_id": 0}}},
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

func (r *blogRepository) FindAll(ctx context.Context, page int64, pageSize int64) ([]models.BlogWithAuthor, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	skip := (page - 1) * pageSize

	pipeline := mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: pageSize}},
	}
	pipeline = append(pipeline, blogAuthorLookup...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []models.BlogWithAuthor
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	if blogs == nil {
		blogs = []models.BlogWithAuthor{}
	}
	fmt.Println("Blogs found:", blogs)

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return blogs, int(totalCount), nil
}

func (r *blogRepository) FindOne(ctx context.Context, id primitive.ObjectID) (*models.BlogWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
	}
	pipeline = append(pipeline, blogAuthorLookup...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.BlogWithAuthor
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	fmt.Printf("Blog found: %+v\n", results[0])

	return &results[0], nil
}

func (r *blogRepository) UpdateOne(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.BlogWithAuthor, error) {
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

	return r.FindOne(ctx, id)
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

func (r *blogRepository) GetBlogsByAuthor(ctx context.Context, authorId primitive.ObjectID, page int64, pageSize int64) ([]models.BlogWithAuthor, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	skip := (page - 1) * pageSize

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"author_id": authorId}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: pageSize}},
	}
	pipeline = append(pipeline, blogAuthorLookup...)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []models.BlogWithAuthor
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	if blogs == nil {
		blogs = []models.BlogWithAuthor{}
	}
	fmt.Printf("Blogs found for author %s: %+v\n", authorId.Hex(), blogs)

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{"author_id": authorId})
	if err != nil {
		return nil, 0, err
	}

	return blogs, int(totalCount), nil
}