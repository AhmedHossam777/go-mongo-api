package repository

import (
	"context"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type refreshTokenRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

type RefreshTokenRepository interface {
	Create(
		ctx context.Context, token *models.RefreshToken,
	) error
	FindActiveTokens(ctx context.Context) (
		[]*models.RefreshToken, error,
	)

	FindActiveTokensByUserID(
		ctx context.Context, userID primitive.ObjectID,
	) ([]*models.RefreshToken, error)

	RevokeToken(
		ctx context.Context, tokenId primitive.ObjectID,
	) error

	RevokeAllUserTokens(
		ctx context.Context, userID primitive.ObjectID,
	) (int64, error)

	DeleteExpiredTokens(ctx context.Context) (
		int64, error,
	)
}

func rewRefreshTokenRepository(db *mongo.Database) RefreshTokenRepository {
	return &refreshTokenRepository{
		collection: db.Collection("refresh_token"),
		timeout:    10 * time.Second,
	}
}

func (r *refreshTokenRepository) Create(
	ctx context.Context, token *models.RefreshToken,
) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *refreshTokenRepository) FindActiveTokens(ctx context.Context) (
	[]*models.RefreshToken, error,
) {
	cursor, err := r.collection.Find(
		ctx, bson.M{
			"revoked":    false,
			"expires_at": bson.M{"$gt": time.Now()},
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var refreshTokens []*models.RefreshToken

	err = cursor.All(ctx, &refreshTokens)
	if err != nil {
		return nil, err
	}

	return refreshTokens, nil
}

func (r *refreshTokenRepository) FindActiveTokensByUserID(
	ctx context.Context, userID primitive.ObjectID,
) ([]*models.RefreshToken, error) {
	cursor, err := r.collection.Find(
		ctx, bson.M{
			"user_id":    userID,
			"revoked":    false,
			"expires_at": bson.M{"$gt": time.Now()},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []*models.RefreshToken
	if err = cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *refreshTokenRepository) RevokeToken(
	ctx context.Context, tokenId primitive.ObjectID,
) error {
	now := time.Now()

	_, err := r.collection.UpdateOne(
		ctx, bson.M{
			"_id": tokenId,
		},
		bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}},
	)

	return err
}

func (r *refreshTokenRepository) RevokeAllUserTokens(
	ctx context.Context, userID primitive.ObjectID,
) (int64, error) {
	now := time.Now()
	result, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "revoked": false},
		bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}},
	)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *refreshTokenRepository) DeleteExpiredTokens(ctx context.Context) (
	int64, error,
) {
	result, err := r.collection.DeleteMany(
		ctx, bson.M{
			"expires_at": bson.M{"$lt": time.Now()},
		},
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
