package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content   string             `json:"content" bson:"content"`
	AuthorId  primitive.ObjectID `json:"author_id" bson:"author_id"`
	BlogId    primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CommentWithAuthor struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content   string             `json:"content" bson:"content"`
	Author    UserResponse       `json:"author" bson:"author"`
	BlogId    primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
