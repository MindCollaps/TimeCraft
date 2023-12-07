package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeTable struct {
	ID   primitive.ObjectID   `json:"id" bson:"_id"`
	Name string               `json:"name" bson:"name"`
	Days []primitive.ObjectID `json:"days" bson:"days"`
}
