package services

import (
	"context"
	"errors"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrBlogNotFound  = errors.New("blog not found")
	ErrInvalidBlogID = errors.New("invalid blog ID")
)

type BlogService interface {
	CreateBlog(ctx context.Context, blogDto *dto.CreateBlogDto, authorId primitive.ObjectID) (*models.Blog, error)
	GetAllBlogs(ctx context.Context, page, pageSize int64) ([]models.Blog, int, error)
	GetBlogByID(ctx context.Context, id string) (*models.Blog, error)
	UpdateBlog(ctx context.Context, id string, updateBlogDto *dto.UpdateBlogDto) (*models.Blog, error)
	DeleteBlog(ctx context.Context, id string) error
}

type blogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) CreateBlog(ctx context.Context, blogDto *dto.CreateBlogDto, authorId primitive.ObjectID) (*models.Blog, error) {
	blog := &models.Blog{
		Title:    blogDto.Title,
		Content:  blogDto.Content,
		AuthorId: authorId,
	}
	return s.repo.Create(ctx, blog)
}

func (s *blogService) GetAllBlogs(ctx context.Context, page, pageSize int64) ([]models.Blog, int, error) {
	return s.repo.FindAll(ctx, page, pageSize)
}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*models.Blog, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidBlogID
	}

	blog, err := s.repo.FindOne(ctx, objectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrBlogNotFound
	}
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *blogService) UpdateBlog(ctx context.Context, id string, updateBlogDto *dto.UpdateBlogDto) (*models.Blog, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidBlogID
	}

	update := bson.M{"updated_at": time.Now()}
	if updateBlogDto.Title != nil {
		update["title"] = *updateBlogDto.Title
	}
	if updateBlogDto.Content != nil {
		update["content"] = *updateBlogDto.Content
	}

	updatedBlog, err := s.repo.UpdateOne(ctx, objectId, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrBlogNotFound
	}
	if err != nil {
		return nil, err
	}

	return updatedBlog, nil
}

func (s *blogService) DeleteBlog(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidBlogID
	}

	err = s.repo.DeleteOne(ctx, objectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrBlogNotFound
	}
	return err
}
