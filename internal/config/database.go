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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	fmt.Println("✅ Database connected successfully!")
	db := client.Database(cfg.DBName)
	err = repository.InitializeIndexes(db)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize indexes: %w", err)
	}

	return db, nil
}
