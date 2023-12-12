package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomConfig struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id"`
	Name           string               `json:"name" bson:"name"`
	RoomNr         string               `json:"roomNr" bson:"roomNr"` // All room number have 2 digits --> 1.0, 3.1, ....
	Capacity       int                  `json:"capacity" bson:"capacity"`
	ExamCapacity   int                  `json:"examCapacity" bson:"examCapacity"`
	Blocks         []primitive.ObjectID `json:"blocks" bson:"blocks"`
	Specialisation []primitive.ObjectID `json:"specialisation" bson:"specialisation"`
}

type RoomConfigStruct struct {
	ID             string               `json:"id" bson:"_id"`
	Name           string               `json:"name" bson:"name"`
	RoomNr         string               `json:"roomNr" bson:"roomNr"` // All room number have 2 digits --> 1.0, 3.1, ....
	Capacity       int                  `json:"capacity" bson:"capacity"`
	ExamCapacity   int                  `json:"examCapacity" bson:"examCapacity"`
	Blocks         []RoomConfig         `json:"blocks" bson:"blocks"`
	Specialisation []RoomSpecialisation `json:"specialisation" bson:"specialisation"`
}
