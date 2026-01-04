package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI string
	DBName   string
	Port     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in environment variables")
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "courseDB"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return &Config{
		MongoURI: mongoURI,
		DBName:   dbName,
		Port:     port,
	}
}

func ConnectDB(cfg *Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).

		// Good balance for typical web API
		SetMaxPoolSize(100).  // Handle ~100 concurrent requests
		SetMinPoolSize(10).   // Keep 10 connections ready
		SetMaxConnecting(10). // Max 10 establishing connections

		// Reasonable timeouts
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(30 * time.Second).

		// Idle connection management
		SetMaxConnIdleTime(60 * time.Second).

		// Enable retries for reliability
		SetRetryWrites(true).
		SetRetryReads(true)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	fmt.Println("✓ Connected to MongoDB with optimized connection pool")
	fmt.Printf("  - Max Pool Size: 100\n")
	fmt.Printf("  - Min Pool Size: 10\n")
	fmt.Printf("  - Max Idle Time: 60s\n")

	db := client.Database(cfg.DBName)
	err = repository.InitializeIndexes(db)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize indexes: %w", err)
	}

	return db, nil
}
