package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type LectureGroup struct {
	Name        string             `json:"name" bson:"name"`
	TimeTableID primitive.ObjectID `json:"timeTableId" bson:"timeTableId"`
}
