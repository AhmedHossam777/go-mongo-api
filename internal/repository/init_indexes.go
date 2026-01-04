package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := initUserIndexes(ctx, db)
	if err != nil {
		fmt.Println("failed to initialize user index, " + err.Error())
		return err
	}

	err = initCourseIndexes(ctx, db)
	if err != nil {
		fmt.Println("failed to initialize course index, " + err.Error())
		return err
	}

	fmt.Println("✓ All indexes initialized successfully")
	return nil
}

func initUserIndexes(ctx context.Context, db *mongo.Database) error {
	userCollection := db.Collection("users")
	indexes := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}

	_, err := userCollection.Indexes().CreateOne(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}

func initCourseIndexes(ctx context.Context, db *mongo.Database) error {
	courseCollection := db.Collection("courses")

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{
					Key:   "course_name",
					Value: 1,
				},
			},
			Options: options.Index().SetUnique(true).SetName("course_name_unique"),
		},
		{
			Keys:    bson.D{{Key: "instructor_id", Value: 1}},
			Options: options.Index().SetName("instructor_index"),
		},
	}

	_, err := courseCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	return nil
}
