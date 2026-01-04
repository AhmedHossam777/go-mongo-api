package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"user_id"`
	Token     string             `json:"-" json:"token"`
	ExpiresAt time.Time          `json:"expiresAt" bson:"expires_at"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	Revoked   bool               `json:"revoked" bson:"revoked"`
	RevokedAt *time.Time         `json:"revokedAt,omitempty" bson:"revoked_at,omitempty"`
	UserAgent string             `json:"userAgent" bson:"user_agent"` // Track user browser
	IPAddress string             `json:"ipAddress" bson:"ip_address"` // Track ip address
}
