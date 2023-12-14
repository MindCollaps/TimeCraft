package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

type Permission struct {
	ID                primitive.ObjectID   `json:"id" bson:"_id"`
	Name              string               `json:"name" bson:"name"`
	IsAdmin           bool                 `json:"isAdmin" bson:"isAdmin"`
	Expires           primitive.DateTime   `json:"username" bson:"username"`
	PermissionEditIds []primitive.ObjectID `json:"permissionEditIds" bson:"permissionEditIds"`
}

type PermissionStruct struct {
	ID                string                 `json:"id" bson:"_id"`
	Name              string                 `json:"name" bson:"name"`
	IsAdmin           bool                   `json:"isAdmin" bson:"isAdmin"`
	Expires           string                 `json:"username" bson:"username"`
	PermissionEditIds []EditPermissionStruct `json:"permissionEditIds" bson:"permissionEditIds"`
}

func PermissionToStruct(c *gin.Context, permission Permission) (PermissionStruct, error) {
	permissionStruct := PermissionStruct{
		ID:      permission.ID.Hex(),
		Name:    permission.Name,
		IsAdmin: permission.IsAdmin,
		Expires: permission.Expires.Time().Format(time.DateTime),
	}

	// Convert PermissionEditIds to []EditPermissionStruct
	editPermissions, err := LoadEditPermissions(c, permission.PermissionEditIds)
	if err != nil {
		return PermissionStruct{}, err
	}

	permissionStruct.PermissionEditIds = editPermissions

	return permissionStruct, nil
}

func LoadPermissions(c *gin.Context, permissionIDs []primitive.ObjectID) ([]PermissionStruct, error) {
	var permissionsStruct []PermissionStruct
	for _, permissionID := range permissionIDs {
		var permission Permission
		err := database.MongoDB.Collection("Permission").FindOne(c, bson.M{
			"_id": permissionID,
		}).Decode(&permission)

		if err != nil {
			continue
		}

		permissionStruct, err := PermissionToStruct(c, permission)
		if err != nil {
			continue
		}
		permissionsStruct = append(permissionsStruct, permissionStruct)
	}
	return permissionsStruct, nil
}
