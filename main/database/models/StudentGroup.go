package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentGroup struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id"`
	Name            string               `json:"name" bson:"name"`
	LectureGroupIds []primitive.ObjectID `json:"lectureGroupIds" bson:"lectureGroupIds"`
	TimeTableId     primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
}
