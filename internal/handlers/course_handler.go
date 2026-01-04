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
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseHandler struct {
	service services.CourseService
}

func NewCourseHandler(service services.CourseService) *CourseHandler {
	return &CourseHandler{
		service: service,
	}
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized,
			"user id not found")
		return
	}
	var courseDto dto.CreateCourseDto
	err := json.NewDecoder(r.Body).Decode(&courseDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest,
			"Invalid request body, "+err.Error())
		return
	}
	defer r.Body.Close()

	mongoUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"error while casting userId "+err.Error())
		return
	}

	courseDto.InstructorId = mongoUserId

	validationErrors := helpers.ValidateStruct(courseDto)
	if validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	createdCourse, err := h.service.CreateCourse(ctx, &courseDto)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			RespondWithError(w, http.StatusBadRequest,
				"duplicated course name")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Error creating course")
		return
	}

	RespondWithJSON(w, http.StatusCreated, createdCourse)
}

func (h *CourseHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
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

	courses, totalCount, err := h.service.GetAllCourses(ctx, int64(page),
		int64(pageSize))

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"Error fetching courses")
		return
	}

	hasMore := false
	count := pageSize
	if totalCount > (page * pageSize) {
		hasMore = true
	} else {
		count = pageSize - (page*pageSize - totalCount)
	}

	PaginationResponse(w, http.StatusOK, courses, page, count, int64(totalCount),
		hasMore)
}

func (h *CourseHandler) GetOneCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	course, err := h.service.GetCourseByID(ctx, courseId)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCourseID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}

		if errors.Is(err, services.ErrCourseNotFound) {
			RespondWithError(w, http.StatusNotFound, "Course not found")
			return
		}

		RespondWithError(w, http.StatusInternalServerError,
			"Error while fetching the course")

		return
	}

	RespondWithJSON(w, http.StatusOK, course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	var updatedCourseDto dto.UpdateCourseDto
	err := json.NewDecoder(r.Body).Decode(&updatedCourseDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest,
			"invalid request body "+err.Error())
		return
	}
	defer r.Body.Close()

	validationErr := helpers.ValidateStruct(updatedCourseDto)

	if validationErr != nil {
		RespondWithValidationErrors(w, validationErr)
		return
	}

	updatedCourse, err := h.service.UpdateCourse(ctx, courseId, &updatedCourseDto)

	RespondWithJSON(w, http.StatusOK, updatedCourse)

}

func (h *CourseHandler) DeleteOneCourse(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	err := h.service.DeleteCourse(ctx, courseId)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCourseID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}
		if errors.Is(err, services.ErrCourseNotFound) {
			RespondWithError(w, http.StatusNotFound, "Course not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Error deleting course")
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}

func (h *CourseHandler) Drop(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.service.Drop(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"Error deleting all courses")
		return
	}
	RespondWithJSON(w, http.StatusOK, nil)
}
