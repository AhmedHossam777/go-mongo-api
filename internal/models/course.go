package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // If ID is empty, don't include it in JSON/BSON
	CourseName string             `json:"course_name" bson:"course_name"`
	Price      int                `json:"price" bson:"course_price"`
	Author     primitive.ObjectID `json:"author" bson:"author"`
}

//type Author struct {
//	Fullname string `json:"fullname" bson:"fullname"`
//	Website  string `json:"website" bson:"website"`
//}
