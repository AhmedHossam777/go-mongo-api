package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// @Summary Create a new product
// @Description Create a new product (requires authentication)
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProduct true "Product details"
// @Success 201 {object} github_com_AhmedHossam777_go-mongo_internal_models.Product "Product created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or duplicate product"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var productDto dto.CreateProduct
	err := json.NewDecoder(r.Body).Decode(&productDto)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"Invalid request body, "+err.Error(),
		)
		return
	}
	defer r.Body.Close()

	validationErrors := helpers.ValidateStruct(productDto)
	if validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	createdProduct, err := h.service.CreateProduct(ctx, &productDto)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			RespondWithError(
				w, http.StatusBadRequest,
				"duplicated product code",
			)
			return
		}
		RespondWithError(
			w, http.StatusInternalServerError, "Error creating product",
		)
		return
	}

	RespondWithJSON(w, http.StatusCreated, createdProduct)
}

// @Summary Get all products
// @Description Get a paginated list of all products
// @Tags products
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Paginated list of products"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default to 10
	}

	products, totalCount, err := h.service.GetAllProducts(
		ctx, int64(page),
		int64(pageSize),
	)

	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"Error fetching products",
		)
		return
	}

	hasMore := false
	count := pageSize
	if totalCount > (page * pageSize) {
		hasMore = true
	} else {
		count = pageSize - (page*pageSize - totalCount)
	}

	PaginationResponse(
		w, http.StatusOK, products, page, count, int64(totalCount),
		hasMore,
	)
}

// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.Product "Product details"
// @Failure 400 {object} map[string]string "Bad request - invalid product ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [get]
func (h *ProductHandler) GetOneProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productId := r.PathValue("id")

	product, err := h.service.GetProductByID(ctx, productId)

	if err != nil {
		if errors.Is(err, services.ErrInvalidProductID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		if errors.Is(err, services.ErrProductNotFound) {
			RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}

		RespondWithError(
			w, http.StatusInternalServerError,
			"Error while fetching the product",
		)

		return
	}

	RespondWithJSON(w, http.StatusOK, product)
}

// @Summary Update product
// @Description Update product details by ID (requires authentication)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProduct true "Updated product details"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.Product "Product updated successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [patch]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productId := r.PathValue("id")

	var updatedProductDto dto.UpdateProduct
	err := json.NewDecoder(r.Body).Decode(&updatedProductDto)
	if err != nil {
		RespondWithError(
			w, http.StatusBadRequest,
			"invalid request body "+err.Error(),
		)
		return
	}
	defer r.Body.Close()

	validationErr := helpers.ValidateStruct(updatedProductDto)

	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	updatedProduct, err := h.service.UpdateProduct(
		ctx, productId, &updatedProductDto,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidProductID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}
		if errors.Is(err, services.ErrProductNotFound) {
			RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		RespondWithError(
			w, http.StatusInternalServerError, "Error updating product",
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, updatedProduct)
}

// @Summary Delete product
// @Description Delete a product by ID (requires authentication)
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 "Product deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid product ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteOneProduct(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productId := r.PathValue("id")

	err := h.service.DeleteProduct(ctx, productId)

	if err != nil {
		if errors.Is(err, services.ErrInvalidProductID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}
		if errors.Is(err, services.ErrProductNotFound) {
			RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		RespondWithError(
			w, http.StatusInternalServerError, "Error deleting product",
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}

// @Summary Drop product collection
// @Description Delete all products from the database (Admin only)
// @Tags products
// @Produce json
// @Success 200 "All products deleted successfully"
// @Failure 403 {object} map[string]string "Forbidden - Admin access required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/drop [delete]
func (h *ProductHandler) Drop(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.service.Drop(ctx)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"Error deleting all products",
		)
		return
	}
	RespondWithJSON(w, http.StatusOK, nil)
}
