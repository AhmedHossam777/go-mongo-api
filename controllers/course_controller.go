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

//! some helper functions

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

//* Create a Course

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	// step 1: create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// step 2 : Decode request body, don't forget to close the body
	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		sendErr(w, http.StatusBadRequest, "Invalid request body, "+err.Error())
		return
	}
	defer r.Body.Close()

	// step 3 : validate the course
	if course.CourseName == "" {
		sendErr(w, http.StatusBadRequest, "course_name is required")
		return
	}

	// step 4 : generate new ObjectID
	course.ID = primitive.NewObjectID()

	// step 5 : insert into mongo DB
	collection := getCourseCollection()
	result, err := collection.InsertOne(ctx,
		course) // this result contains the insertID
	if err != nil {
		sendErr(w, http.StatusInternalServerError,
			"Error while creating a course "+err.Error())
		return
	}

	// step 6 : send response
	sendJson(w, http.StatusCreated, map[string]interface{}{
		"message":  "course created successfully",
		"insertID": result.InsertedID,
		"course":   course,
	})
}

// GET all courses

func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	// step 1 : make the context
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

	// Step 3: Handle empty result
	if courses == nil {
		courses = []models.Course{} // Return empty array, not null
	}

	sendJson(w, http.StatusOK, map[string]interface{}{
		"message": "fetching all courses successfully",
		"courses": courses,
	})
}
