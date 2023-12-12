package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeTableDay struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Date        primitive.DateTime   `json:"date" bson:"date"`
	TimeSlotIds []primitive.ObjectID `json:"timeSlotIds" bson:"timeSlotIds"`
	LastUpdated primitive.DateTime   `json:"lastUpdated" bson:"lastUpdated"`
}

type TimeTableDayStruct struct {
	ID          string     `json:"id" bson:"_id"`
	Date        string     `json:"date" bson:"date"`
	TimeSlotIds []TimeSlot `json:"timeSlotIds" bson:"timeSlotIds"`
	LastUpdated string     `json:"lastUpdated" bson:"lastUpdated"`
}
