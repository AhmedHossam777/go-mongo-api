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
	CreateBlog(ctx context.Context, blogDto *dto.CreateBlogDto, authorId primitive.ObjectID) (*models.BlogWithAuthor, error)
	GetAllBlogs(ctx context.Context, page, pageSize int64) ([]models.BlogWithAuthor, int, error)
	GetBlogByID(ctx context.Context, id string) (*models.BlogWithAuthor, error)
	UpdateBlog(ctx context.Context, id string, updateBlogDto *dto.UpdateBlogDto) (*models.BlogWithAuthor, error)
	DeleteBlog(ctx context.Context, id string) error
	AddComment(ctx context.Context, blogId string, authorId primitive.ObjectID, commentDto *dto.AddCommentDto) (*models.CommentWithAuthor, error)
	GetCommentsByBlogId(ctx context.Context, blogId string) ([]models.CommentWithAuthor, error)
	GetBlogsByAuthor(ctx context.Context, authorId primitive.ObjectID, page, pageSize int64) ([]models.BlogWithAuthor, int, error)
}

type blogService struct {
	repo        repository.BlogRepository
	commentRepo repository.CommentRepository
}

func NewBlogService(repo repository.BlogRepository, commentRepo repository.CommentRepository) BlogService {
	return &blogService{repo: repo, commentRepo: commentRepo}
}

func (s *blogService) CreateBlog(ctx context.Context, blogDto *dto.CreateBlogDto, authorId primitive.ObjectID) (*models.BlogWithAuthor, error) {
	blog := &models.Blog{
		Title:    blogDto.Title,
		Content:  blogDto.Content,
		ImageURL: blogDto.ImageURL,
		AuthorId: authorId,
	}
	created, err := s.repo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}
	return s.repo.FindOne(ctx, created.ID)
}

func (s *blogService) GetAllBlogs(ctx context.Context, page, pageSize int64) ([]models.BlogWithAuthor, int, error) {
	return s.repo.FindAll(ctx, page, pageSize)
}

func (s *blogService) GetBlogsByAuthor(ctx context.Context, authorId primitive.ObjectID, page, pageSize int64) ([]models.BlogWithAuthor, int, error) {
	return s.repo.GetBlogsByAuthor(ctx, authorId, page, pageSize)
}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*models.BlogWithAuthor, error) {
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

func (s *blogService) UpdateBlog(ctx context.Context, id string, updateBlogDto *dto.UpdateBlogDto) (*models.BlogWithAuthor, error) {
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
	if err != nil {
		if err.Error() == "blog not found" {
			return nil, ErrBlogNotFound
		}
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

func (s *blogService) AddComment(ctx context.Context, blogId string, authorId primitive.ObjectID, commentDto *dto.AddCommentDto) (*models.CommentWithAuthor, error) {
	blogObjectId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, ErrInvalidBlogID
	}

	_, err = s.repo.FindOne(ctx, blogObjectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrBlogNotFound
	}
	if err != nil {
		return nil, err
	}

	comment := &models.Comment{
		Content:  commentDto.Content,
		AuthorId: authorId,
		BlogId:   blogObjectId,
	}

	return s.commentRepo.Create(ctx, comment)
}

func (s *blogService) GetCommentsByBlogId(ctx context.Context, blogId string) ([]models.CommentWithAuthor, error) {
	blogObjectId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, ErrInvalidBlogID
	}

	_, err = s.repo.FindOne(ctx, blogObjectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrBlogNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.commentRepo.FindByBlogId(ctx, blogObjectId)
}
