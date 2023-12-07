package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LecturersPreferedLecture struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	LectureId  primitive.ObjectID `json:"lectureId" bson:"lectureId"`
	MaxPerWeek int                `json:"maxPerWeek" bson:"maxPerWeek"`
	MaxPerDay  int                `json:"maxPerDay" bson:"maxPerDay"`
	Type       int                `json:"type" bson:"type"`
}
