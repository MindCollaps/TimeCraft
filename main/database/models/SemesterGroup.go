package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SemesterGroup struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Name               string               `json:"name" bson:"name"`
	StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableId        primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds" bson:"specialisationsIds"`
}

type SemesterGroupStruct struct {
	ID                 string         `json:"id" bson:"_id"`
	Name               string         `json:"name" bson:"name"`
	StudentGroupIds    []StudentGroup `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableId        TimeTable      `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []StudentGroup `json:"specialisationsIds" bson:"specialisationsIds"`
}
