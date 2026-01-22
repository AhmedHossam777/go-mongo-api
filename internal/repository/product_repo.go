package repository

import (
	"context"
	"errors"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	FindAll(ctx context.Context, page int64, pageSize int64) (
		[]models.Product, int, error,
	)
	FindOne(ctx context.Context, productId primitive.ObjectID) (
		*models.Product, error,
	)
	UpdateOne(
		ctx context.Context, productId primitive.ObjectID, update bson.M,
	) (
		*models.Product, error,
	)
	DeleteOne(ctx context.Context, productId primitive.ObjectID) error
	Drop(ctx context.Context) error
}

type productRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewProductRepo(db *mongo.Database) ProductRepository {
	return &productRepository{
		collection: db.Collection("products"),
		timeout:    10 * time.Second,
	}
}

func (r *productRepository) Create(
	ctx context.Context, product *models.Product,
) (*models.Product, error) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	product.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) FindAll(
	ctx context.Context, page int64, pageSize int64,
) (
	[]models.Product, int, error,
) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	skip := (page - 1) * pageSize

	findOptions := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}). // Sort by newest first
		SetSkip(skip).
		SetLimit(pageSize)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, 0, err
	}

	// if products table is empty return empty slice not nil
	if products == nil {
		products = []models.Product{}
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return products, int(totalCount), nil
}

func (r *productRepository) FindOne(
	ctx context.Context, id primitive.ObjectID,
) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var product *models.Product

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) UpdateOne(
	ctx context.Context, id primitive.ObjectID, update bson.M,
) (*models.Product, error) {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedProduct models.Product

	err := r.collection.FindOneAndUpdate(
		ctx, filter, update,
		opts,
	).Decode(&updatedProduct)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *productRepository) DeleteOne(
	ctx context.Context, id primitive.ObjectID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	deleteResult, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *productRepository) Drop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	err := r.collection.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}
