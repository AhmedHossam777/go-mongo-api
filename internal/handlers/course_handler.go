package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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

	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request body, "+err.Error())
		return
	}
	defer r.Body.Close()

	createdCourse, err := h.service.CreateCourse(ctx, &course)

	if err != nil {
		if errors.Is(err, services.ErrCourseNameRequired) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "Error creating course")
		return
	}

	SendSuccess(w, http.StatusCreated, "Course created successfully",
		createdCourse)
}

func (h *CourseHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courses, err := h.service.GetAllCourses(ctx)
	if err != nil {
		log.Println(err)
		SendError(w, http.StatusInternalServerError, "Error fetching courses")
		return
	}

	SendSuccess(w, http.StatusOK, "Courses fetched successfully", courses)
}

func (h *CourseHandler) GetOneCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	course, err := h.service.GetCourseByID(ctx, courseId)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCourseID) {
			SendError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}

		if errors.Is(err, services.ErrCourseNotFound) {
			SendError(w, http.StatusNotFound, "Course not found")
			return
		}

		SendError(w, http.StatusInternalServerError,
			"Error while fetching the course")

		return
	}

	SendSuccess(w, http.StatusOK, "Course fetched successfully", course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	var updatedData models.Course
	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		SendError(w, http.StatusBadRequest, "invalid request body "+err.Error())
		return
	}
	defer r.Body.Close()

	updatedCourse, err := h.service.UpdateCourse(ctx, courseId, &updatedData)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCourseID) {
			SendError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}
		if errors.Is(err, services.ErrCourseNotFound) {
			SendError(w, http.StatusNotFound, "Course not found")
			return
		}
		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			SendError(w, http.StatusBadRequest, "Not enough fields to update")
			return
		}

		SendError(w, http.StatusInternalServerError, "Error while updating course")
	}

	SendSuccess(w, http.StatusOK, "course updated successfully", updatedCourse)

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
			SendError(w, http.StatusBadRequest, "Invalid course ID")
			return
		}
		if errors.Is(err, services.ErrCourseNotFound) {
			SendError(w, http.StatusNotFound, "Course not found")
			return
		}
		SendError(w, http.StatusInternalServerError, "Error deleting course")
		return
	}

	SendSuccess(w, http.StatusOK, "Course deleted successfully", nil)
}
