package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SemesterGroup struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Name               string               `json:"name" bson:"name"`
	StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableID        primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds" bson:"specialisationsIds"`
}
