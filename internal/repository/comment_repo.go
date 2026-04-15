package repository

import (
	"context"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) (*models.CommentWithAuthor, error)
	FindByBlogId(ctx context.Context, blogId primitive.ObjectID) ([]models.CommentWithAuthor, error)
}

type commentRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewCommentRepo(db *mongo.Database) CommentRepository {
	return &commentRepository{
		collection: db.Collection("comments"),
		timeout:    10 * time.Second,
	}
}

var authorLookupStage = bson.D{
	{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "author_id",
		"foreignField": "_id",
		"as":           "authorArr",
	}},
}

var authorUnwindStage = bson.D{
	{Key: "$unwind", Value: "$authorArr"},
}

var authorProjectStage = bson.D{
	{Key: "$addFields", Value: bson.M{
		"author": bson.M{
			"id":        "$authorArr._id",
			"name":      "$authorArr.name",
			"email":     "$authorArr.email",
			"role":      "$authorArr.role",
			"createdAt": "$authorArr.created_at",
		},
	}},
}

var removeAuthorArrStage = bson.D{
	{Key: "$project", Value: bson.M{"authorArr": 0}},
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) (*models.CommentWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": comment.ID}}},
		authorLookupStage,
		authorUnwindStage,
		authorProjectStage,
		removeAuthorArrStage,
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.CommentWithAuthor
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

func (r *commentRepository) FindByBlogId(ctx context.Context, blogId primitive.ObjectID) ([]models.CommentWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"blog_id": blogId}}},
		authorLookupStage,
		authorUnwindStage,
		authorProjectStage,
		removeAuthorArrStage,
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.CommentWithAuthor
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	if comments == nil {
		comments = []models.CommentWithAuthor{}
	}

	return comments, nil
}
