package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lecturer struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	FirstName    string               `json:"firstName" bson:"firstName"`
	SureName     string               `json:"lastName" bson:"lastName"`
	Title        string               `json:"title" bson:"title"`
	CanHold      []primitive.ObjectID `json:"canHold" bson:"canHold"`
	ContactEmail string               `json:"contactEmail" bson:"contactEmail"`
	ContactPhone string               `json:"contactPhone" bson:"contactPhone"`
	MaxWeekHours int                  `json:"maxWeekHours" bson:"maxWeekHours"`
}

type LecturerStruct struct {
	ID           string    `json:"id" bson:"_id"`
	FirstName    string    `json:"firstName" bson:"firstName"`
	SureName     string    `json:"lastName" bson:"lastName"`
	Title        string    `json:"title" bson:"title"`
	CanHold      []Lecture `json:"canHold" bson:"canHold"`
	ContactEmail string    `json:"contactEmail" bson:"contactEmail"`
	ContactPhone string    `json:"contactPhone" bson:"contactPhone"`
	MaxWeekHours int       `json:"maxWeekHours" bson:"maxWeekHours"`
}
