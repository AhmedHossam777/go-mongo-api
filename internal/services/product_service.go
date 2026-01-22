package services

import (
	"context"
	"errors"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInvalidProductID = errors.New("invalid product ID")
)

type ProductService interface {
	CreateProduct(ctx context.Context, productDto *dto.CreateProduct) (
		*models.Product, error,
	)
	GetAllProducts(ctx context.Context, page, pageSize int64) (
		[]models.Product, int, error,
	)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	UpdateProduct(
		ctx context.Context, id string, updateProductDto *dto.UpdateProduct,
	) (*models.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	Drop(ctx context.Context) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(
	ctx context.Context, productDto *dto.CreateProduct,
) (*models.Product, error) {

	product := &models.Product{
		Title:       productDto.Title,
		Price:       productDto.Price,
		Description: productDto.Description,
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetAllProducts(
	ctx context.Context, page, pageSize int64,
) (
	[]models.Product, int, error,
) {
	return s.repo.FindAll(ctx, page, pageSize)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (
	*models.Product, error,
) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidProductID
	}

	product, err := s.repo.FindOne(ctx, objectId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(
	ctx context.Context, id string, updateProductDto *dto.UpdateProduct,
) (*models.Product, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidProductID
	}

	var update = bson.M{}
	if updateProductDto.Title != "" {
		update["title"] = updateProductDto.Title
	}

	if updateProductDto.Price > 0 {
		update["price"] = updateProductDto.Price
	}
	if updateProductDto.Description != "" {
		update["description"] = updateProductDto.Description
	}

	updateProduct, err := s.repo.UpdateOne(ctx, objectId, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return updateProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidProductID
	}

	err = s.repo.DeleteOne(ctx, objectId)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrProductNotFound
	}
	return nil
}

func (s *productService) Drop(ctx context.Context) error {
	err := s.repo.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}
