package helpers

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(
	userId primitive.ObjectID, email string, role string,
) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	expirationHours := 24
	hours := os.Getenv("JWT_EXPIRATION_HOURS")
	if hours != "" {
		h, err := strconv.Atoi(hours)
		if err == nil {
			expirationHours = h
		}
	}

	claims := JWTClaims{
		UserId: userId.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "courses-api",
			Subject:   userId.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// get user id from token claims
func GetUserIdFromToken(claims *JWTClaims) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(claims.UserId)
}
