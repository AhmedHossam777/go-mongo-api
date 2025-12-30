package repository

import (
	"context"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetOneUser(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	UpdateOneUser(ctx context.Context, id primitive.ObjectID, update bson.M) (
		*models.User, error,
	)
	DeleteOneUser(ctx context.Context, id primitive.ObjectID) error
}

func NewUserRepo(db *mongo.Database) UserRepository {
	return &userRepo{
		collection: db.Collection("users"),
		timeout:    10 * time.Second,
	}
}

func (r *userRepo) CreateUser(
	ctx context.Context, user *models.User,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	user.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []models.User

	err = cursor.All(ctx, &users)

	if err != nil {
		return nil, err
	}

	if users == nil {
		users = []models.User{}
	}

	return users, nil
}

func (r *userRepo) GetOneUser(
	ctx context.Context, id primitive.ObjectID,
) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user *models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) UpdateOneUser(
	ctx context.Context, id primitive.ObjectID, update bson.M,
) (
	*models.User, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	updateResult, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)

	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	fetchCtx, cancelFetchCtx := context.WithTimeout(ctx, r.timeout)
	defer cancelFetchCtx()
	return r.GetOneUser(fetchCtx, id)
}

func (r *userRepo) DeleteOneUser(
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
