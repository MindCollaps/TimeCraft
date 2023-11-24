package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SemesterGroup struct {
	Name               string               `json:"name" bson:"name"`
	StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableID        primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds" bson:"specialisationsIds"`
}
