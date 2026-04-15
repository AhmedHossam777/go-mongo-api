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

const maxImageSize = 10 << 20 // 10 MB

type BlogHandler struct {
	service   services.BlogService
	s3Service services.S3Service
}

func NewBlogHandler(service services.BlogService, s3Service services.S3Service) *BlogHandler {
	return &BlogHandler{service: service, s3Service: s3Service}
}

// @Summary Create a new blog
// @Description Create a new blog post (requires authentication). Send as multipart/form-data with fields: title, content, and optional image file.
// @Tags blogs
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Blog title"
// @Param content formData string true "Blog content"
// @Param image formData file false "Blog cover image"
// @Success 201 {object} github_com_AhmedHossam777_go-mongo_internal_models.Blog "Blog created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs [post]
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "user id not found")
		return
	}

	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		RespondWithError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	blogDto := dto.CreateBlogDto{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	validationErrors := helpers.ValidateStruct(blogDto)
	if validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageURL, err := h.s3Service.UploadFile(ctx, file, header, "blogs")
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "failed to upload image: "+err.Error())
			return
		}
		blogDto.ImageURL = imageURL
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

// @Summary Get blogs by author
// @Description Get a paginated list of blog posts created by the logged-in user (requires authentication)
// @Tags blogs
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} map[string]interface{} "Paginated list of blogs"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/my [get]	
func (h *BlogHandler) GetBlogsByAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "user id not found")
		return
	}

	authorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error while casting userId "+err.Error())
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	blogs, totalCount, err := h.service.GetBlogsByAuthor(ctx, authorId, int64(page), int64(pageSize))
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

// @Summary Add a comment to a blog post
// @Description Add a comment to a blog post written by the logged-in user
// @Tags blogs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param request body dto.AddCommentDto true "Comment content"
// @Success 201 {object} github_com_AhmedHossam777_go-mongo_internal_models.CommentWithAuthor "Comment created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Blog not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/{id}/comments [post]
func (h *BlogHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "user id not found")
		return
	}

	authorId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error while casting userId "+err.Error())
		return
	}

	blogId := r.PathValue("id")

	var commentDto dto.AddCommentDto
	if err := json.NewDecoder(r.Body).Decode(&commentDto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if validationErrors := helpers.ValidateStruct(commentDto); validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	comment, err := h.service.AddComment(ctx, blogId, authorId, &commentDto)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBlogID) {
			RespondWithError(w, http.StatusBadRequest, "invalid blog ID")
			return
		}
		if errors.Is(err, services.ErrBlogNotFound) {
			RespondWithError(w, http.StatusNotFound, "blog not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "error adding comment")
		return
	}

	RespondWithJSON(w, http.StatusCreated, comment)
}

// @Summary Get comments for a blog post
// @Description Get all comments for a specific blog post
// @Tags blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {array} github_com_AhmedHossam777_go-mongo_internal_models.CommentWithAuthor "List of comments"
// @Failure 400 {object} map[string]string "Bad request - invalid blog ID"
// @Failure 404 {object} map[string]string "Blog not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /blogs/{id}/comments [get]
func (h *BlogHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogId := r.PathValue("id")

	comments, err := h.service.GetCommentsByBlogId(ctx, blogId)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBlogID) {
			RespondWithError(w, http.StatusBadRequest, "invalid blog ID")
			return
		}
		if errors.Is(err, services.ErrBlogNotFound) {
			RespondWithError(w, http.StatusNotFound, "blog not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "error fetching comments")
		return
	}

	RespondWithJSON(w, http.StatusOK, comments)
}
