package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomSpecialisation struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type RoomSpecialisationStruct struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
