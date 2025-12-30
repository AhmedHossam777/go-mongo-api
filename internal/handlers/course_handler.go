package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/dto"
	"github.com/AhmedHossam777/go-mongo/internal/helpers"
	"github.com/AhmedHossam777/go-mongo/internal/models"
	"github.com/AhmedHossam777/go-mongo/internal/services"
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

	var courseDto dto.CreateCourseDto
	err := json.NewDecoder(r.Body).Decode(&courseDto)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest,
			"Invalid request body, "+err.Error())
		return
	}

	defer r.Body.Close()

	validationErrors := helpers.ValidateStruct(courseDto)
	if validationErrors != nil {
		RespondWithValidationErrors(w, validationErrors)
		return
	}

	createdCourse, err := h.service.CreateCourse(ctx, &courseDto)

	if err != nil {
		if errors.Is(err, services.ErrCourseNameRequired) {
			RespondWithError(w, http.StatusBadRequest, err.Error())
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

	courses, err := h.service.GetAllCourses(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError,
			"Error fetching courses")
		return
	}

	RespondWithJSON(w, http.StatusOK, courses)
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

	var updatedData models.Course
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest,
			"invalid request body "+err.Error())
		return
	}
	defer r.Body.Close()

	updatedCourse, err := h.service.UpdateCourse(ctx, courseId, &updatedData)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCourseID) {
			RespondWithError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}
		if errors.Is(err, services.ErrCourseNotFound) {
			RespondWithError(w, http.StatusNotFound, "Course not found")
			return
		}
		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			RespondWithError(w, http.StatusBadRequest, "Not enough fields to update")
			return
		}

		RespondWithError(w, http.StatusInternalServerError,
			"Error while updating course")
	}

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
