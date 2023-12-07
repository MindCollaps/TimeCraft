package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type LectureGroup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	TimeTableId primitive.ObjectID `json:"timeTableId" bson:"timeTableId"`
}
