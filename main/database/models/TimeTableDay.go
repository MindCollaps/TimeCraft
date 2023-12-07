package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeTableDay struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Date        primitive.DateTime   `json:"date" bson:"date"`
	TimeSlotIDs []primitive.ObjectID `json:"timeSlotIDs" bson:"timeSlotIDs"`
	LastUpdated primitive.DateTime   `json:"lastUpdated" bson:"lastUpdated"`
}
