package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	ID                primitive.ObjectID   `json:"id" bson:"_id"`
	Name              string               `json:"name" bson:"name"`
	IsAdmin           bool                 `json:"isAdmin" bson:"isAdmin"`
	Expires           primitive.DateTime   `json:"username" bson:"username"`
	PermissionEditIds []primitive.ObjectID `json:"permissionEditIds" bson:"permissionEditIds"`
}
