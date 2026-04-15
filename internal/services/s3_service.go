package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error)
}

type s3Service struct {
	client    *s3.Client
	bucket    string
	baseURL   string
}

func NewS3Service() (S3Service, error) {
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("AWS_S3_BUCKET")
	baseURL := os.Getenv("AWS_S3_BASE_URL")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &s3Service{
		client:  client,
		bucket:  bucket,
		baseURL: strings.TrimRight(baseURL, "/"),
	}, nil
}

func (s *s3Service) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}
