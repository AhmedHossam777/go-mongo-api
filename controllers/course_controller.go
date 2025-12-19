package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AhmedHossam777/go-mongo/config"
	"github.com/AhmedHossam777/go-mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCourseCollection() *mongo.Collection {
	return config.DB.Collection("courses")
}

func sendJson(
	w http.ResponseWriter, status int, data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		sendErr(w, 500, "error while encoding json response")
	}
}

func sendErr(
	w http.ResponseWriter, status int, message string,
) {
	sendJson(w, status, map[string]string{"message": message})
}

// * Create a Course

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		sendErr(w, http.StatusBadRequest, "Invalid request body, "+err.Error())
		return
	}
	defer r.Body.Close()

	if course.CourseName == "" {
		sendErr(w, http.StatusBadRequest, "course_name is required")
		return
	}

	course.ID = primitive.NewObjectID()

	collection := getCourseCollection()
	result, err := collection.InsertOne(ctx,
		course) // this result contains the insertID
	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error while creating a course "+err.Error())
		return
	}

	sendJson(w, http.StatusCreated, map[string]interface{}{
		"message":  "course created successfully",
		"insertID": result.InsertedID,
		"course":   course,
	})
}

func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := getCourseCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"error while getting all courses "+err.Error())
		return
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	err = cursor.All(ctx, &courses)
	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error decoding courses: "+err.Error())
		return
	}

	if courses == nil {
		courses = []models.Course{}
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"message": "fetching all courses successfully",
		"courses": courses,
	})
}

func GetOneCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(courseId)

	if err != nil {
		sendErr(w, http.StatusBadRequest,
			"Invalid course id please add a valid mongo id")
		return
	}

	collection := getCourseCollection()
	var course models.Course

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&course)

	if err == mongo.ErrNoDocuments {
		sendErr(w, http.StatusNotFound, "Course not found")
		return
	}

	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error fetching course: "+err.Error())
		return
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"message": "course found successfully",
		"course":  course,
	})
}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		sendErr(w, http.StatusBadRequest, "invalid mongo id")
		return
	}

	var updatedData models.Course
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		sendErr(w, http.StatusBadRequest, "invalid request body "+err.Error())
		return
	}
	defer r.Body.Close()

	update := bson.M{
		"$set": bson.M{},
	}

	setFields := update["$set"].(bson.M)
	if updatedData.CourseName != "" {
		setFields["course_name"] = updatedData.CourseName
	}
	if updatedData.Price != 0 {
		setFields["course_price"] = updatedData.Price
	}
	if updatedData.Author != nil {
		setFields["author"] = updatedData.Author
	}

	if len(setFields) == 0 {
		sendErr(w, http.StatusBadRequest, "no data to update")
		return
	}

	collection := getCourseCollection()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)

	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error updating course: "+err.Error())
		return
	}

	if result.MatchedCount == 0 {
		sendErr(w, http.StatusNotFound, "course not found")
		return
	}

	var updatedCourse models.Course
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedCourse)

	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error fetching updated course: "+err.Error())
		return
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"message": "course update successfully",
		"course":  updatedCourse,
	})
}

func DeleteOneCourse(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(courseId)

	if err != nil {
		sendErr(w, http.StatusBadRequest,
			"Invalid course id please add a valid mongo id")
		return
	}

	collection := getCourseCollection()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})

	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error deleting course: "+err.Error())
		return
	}

	if result.DeletedCount == 0 {
		sendErr(w, http.StatusNotFound, "course not found")
		return
	}

	sendJson(w, http.StatusOK, map[string]string{
		"message": "course deleted successfully",
	})
}
