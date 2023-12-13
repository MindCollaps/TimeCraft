package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type EditPermission struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	DbName      string             `json:"dbName" bson:"dbName"`
	DbElementId primitive.ObjectID `json:"dbElementId" bson:"dbElementId"`
}

type EditPermissionStruct struct {
	ID          string `json:"id" bson:"_id"`
	DbName      string `json:"dbName" bson:"dbName"`
	DbElementId string `json:"dbElementId" bson:"dbElementId"`
}

func EditPermissionToStruct(editPermission EditPermission) EditPermissionStruct {
	editPermissionStruct := EditPermissionStruct{
		ID:          editPermission.ID.Hex(),
		DbName:      editPermission.DbName,
		DbElementId: editPermission.DbElementId.Hex(),
	}
	return editPermissionStruct
}

func LoadEditPermissions(c *gin.Context, editPermissionIDs []primitive.ObjectID) ([]EditPermissionStruct, error) {
	var editPermissionsStruct []EditPermissionStruct
	for _, editPermissionID := range editPermissionIDs {
		var editPermission EditPermission
		err := database.MongoDB.Collection("EditPermissions").FindOne(c, bson.M{
			"_id": editPermissionID,
		}).Decode(&editPermission)

		if err != nil {
			continue
		}

		editPermissionStruct := EditPermissionToStruct(editPermission)
		editPermissionsStruct = append(editPermissionsStruct, editPermissionStruct)
	}
	return editPermissionsStruct, nil
}
