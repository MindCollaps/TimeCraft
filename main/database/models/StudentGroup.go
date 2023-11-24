package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentGroup struct {
	Name            string               `json:"name" bson:"name"`
	LectureGroupIds []primitive.ObjectID `json:"lectureGroupIds" bson:"lectureGroupIds"`
	TimeTableID     primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
}
