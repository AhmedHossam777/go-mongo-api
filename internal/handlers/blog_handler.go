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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogHandler struct {
	service services.BlogService
}

func NewBlogHandler(service services.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

// @Summary Create a new blog
// @Description Create a new blog post (requires authentication)
// @Tags blogs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateBlogDto true "Blog details"
// @Success 201 {object} github_com_AhmedHossam777_go-mongo_internal_models.Blog "Blog created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs [post]
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "user id not found")
		return
	}

	var blogDto dto.CreateBlogDto
	err := json.NewDecoder(r.Body).Decode(&blogDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body, "+err.Error())
		return
	}
	defer r.Body.Close()

	validationErrors := helpers.ValidateStruct(blogDto)
	if validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	authorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error while casting userId "+err.Error())
		return
	}

	createdBlog, err := h.service.CreateBlog(ctx, &blogDto, authorId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error creating blog")
		return
	}

	RespondWithJSON(w, http.StatusCreated, createdBlog)
}

// @Summary Get all blogs
// @Description Get a paginated list of all blog posts
// @Tags blogs
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Paginated list of blogs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs [get]
func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	blogs, totalCount, err := h.service.GetAllBlogs(ctx, int64(page), int64(pageSize))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error fetching blogs")
		return
	}

	hasMore := false
	count := pageSize
	if totalCount > (page * pageSize) {
		hasMore = true
	} else {
		count = pageSize - (page*pageSize - totalCount)
	}

	PaginationResponse(w, http.StatusOK, blogs, page, count, int64(totalCount), hasMore)
}

// @Summary Get blog by ID
// @Description Get a specific blog post by its ID
// @Tags blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.Blog "Blog details"
// @Failure 400 {object} map[string]string "Bad request - invalid blog ID"
// @Failure 404 {object} map[string]string "Blog not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/{id} [get]
func (h *BlogHandler) GetOneBlog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := r.PathValue("id")

	blog, err := h.service.GetBlogByID(ctx, blogId)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBlogID) {
			RespondWithError(w, http.StatusBadRequest, "invalid blog ID")
			return
		}
		if errors.Is(err, services.ErrBlogNotFound) {
			RespondWithError(w, http.StatusNotFound, "blog not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "error while fetching the blog")
		return
	}

	RespondWithJSON(w, http.StatusOK, blog)
}

// @Summary Update blog
// @Description Update blog post details by ID (requires authentication)
// @Tags blogs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param request body dto.UpdateBlogDto true "Updated blog details"
// @Success 200 {object} github_com_AhmedHossam777_go-mongo_internal_models.Blog "Blog updated successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Blog not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/{id} [patch]
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := r.PathValue("id")

	var updateBlogDto dto.UpdateBlogDto
	err := json.NewDecoder(r.Body).Decode(&updateBlogDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body "+err.Error())
		return
	}
	defer r.Body.Close()

	validationErr := helpers.ValidateStruct(updateBlogDto)
	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	updatedBlog, err := h.service.UpdateBlog(ctx, blogId, &updateBlogDto)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBlogID) {
			RespondWithError(w, http.StatusBadRequest, "invalid blog ID")
			return
		}
		if errors.Is(err, services.ErrBlogNotFound) {
			RespondWithError(w, http.StatusNotFound, "blog not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "error updating blog")
		return
	}

	RespondWithJSON(w, http.StatusOK, updatedBlog)
}

// @Summary Delete blog
// @Description Delete a blog post by ID (requires authentication)
// @Tags blogs
// @Security BearerAuth
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 "Blog deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid blog ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Blog not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/{id} [delete]
func (h *BlogHandler) DeleteOneBlog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := r.PathValue("id")

	err := h.service.DeleteBlog(ctx, blogId)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBlogID) {
			RespondWithError(w, http.StatusBadRequest, "invalid blog ID")
			return
		}
		if errors.Is(err, services.ErrBlogNotFound) {
			RespondWithError(w, http.StatusNotFound, "blog not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "error deleting blog")
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
