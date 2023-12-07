package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EditPermission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	DbName      string             `json:"dbName" bson:"dbName"`
	DbElementId primitive.ObjectID `json:"dbElementId" bson:"dbElementId"`
}
