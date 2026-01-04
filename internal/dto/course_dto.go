package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateCourseDto struct {
	CourseName   string             `json:"course_name" validate:"required,min=2,max=100"`
	Price        int                `json:"price" validate:"required,gt=0"`
	InstructorId primitive.ObjectID `json:"instructor_id" validate:"required"`
}
type UpdateCourseDto struct {
	CourseName *string `json:"course_name" validate:"omitempty,min=2,max=100"`
	Price      *int    `json:"price" validate:"omitempty,gte=0"`
}
