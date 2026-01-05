package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	expirationMinutes := 15
	mins := os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES")
	if mins != "" {
		m, err := strconv.Atoi(mins)
		if err == nil {
			expirationMinutes = m
		}
	}

	claims := JWTClaims{
		UserId: userId.Hex(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "courses-api",
			Subject:   userId.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
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

func GenerateRefreshToken() (string, error) {
	// Generate 32 random bytes (256 bits of entropy)
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashRefreshToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func ValidateRefreshToken(token, hash string) bool {
	token = strings.TrimSpace(token)
	hash = strings.TrimSpace(hash)

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))

	return err == nil
}

func GetRefreshTokenExpiry() time.Time {
	days := 7

	d := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS")

	if d != "" {
		parsed, err := strconv.Atoi(d)
		if err != nil {
			days = parsed
		}
	}
	return time.Now().Add(time.Duration(days) * 24 * time.Hour)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(
		tokenString, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
