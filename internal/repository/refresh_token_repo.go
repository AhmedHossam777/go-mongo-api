package repository

import (
	"context"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		collection: db.Collection("refresh_token"),
		timeout:    10 * time.Second,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context, token *models.RefreshToken,
) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *RefreshTokenRepository) FindActiveTokens(ctx context.Context) (
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

func (r *RefreshTokenRepository) FindActiveTokensByUserID(
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

func (r *RefreshTokenRepository) RevokeToken(
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

func (r *RefreshTokenRepository) RevokeAllUserTokens(
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

func (r *RefreshTokenRepository) DeleteExpiredTokens(ctx context.Context) (
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
