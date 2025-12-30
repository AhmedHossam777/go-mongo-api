package dto

type CreateCourseDto struct {
	CourseName string `json:"course_name" validate:"required,min=2,max=100"`
	Price      int    `json:"price" validate:"required,min=1"`
}
type UpdateCourseDto struct {
	CourseName string `json:"course_name" validate:"omitempty,min=2,max=100"`
	Price      int    `json:"price" validate:"omitempty,min=0"`
}
