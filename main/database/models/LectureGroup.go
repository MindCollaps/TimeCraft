package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type LectureGroup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	TimeTableId primitive.ObjectID `json:"timeTableId" bson:"timeTableId"`
}

type LectureGroupStruct struct {
	ID          string          `json:"id" bson:"_id"`
	Name        string          `json:"name" bson:"name"`
	TimeTableId TimeTableStruct `json:"timeTableId" bson:"timeTableId"`
}

func LectureGroupToStruct(c *gin.Context, lectureGroup LectureGroup) (LectureGroupStruct, error) {
	lectureGroupStruct := LectureGroupStruct{
		ID:   lectureGroup.ID.Hex(),
		Name: lectureGroup.Name,
	}

	// Convert TimeTableId to TimeTableStruct
	timeTables, err := LoadTimeTables(c, []primitive.ObjectID{lectureGroup.TimeTableId})
	if err != nil {
		return LectureGroupStruct{}, err
	}

	lectureGroupStruct.TimeTableId = timeTables[0]

	return lectureGroupStruct, nil
}

func LoadLectureGroups(c *gin.Context, lectureGroupIDs []primitive.ObjectID) ([]LectureGroupStruct, error) {
	var lectureGroupsStruct []LectureGroupStruct
	for _, lectureGroupID := range lectureGroupIDs {
		var lectureGroup LectureGroup
		err := database.MongoDB.Collection("LectureGroup").FindOne(c, bson.M{
			"_id": lectureGroupID,
		}).Decode(&lectureGroup)

		if err != nil {
			continue
		}

		lectureGroupStruct, err := LectureGroupToStruct(c, lectureGroup)
		if err != nil {
			continue
		}
		lectureGroupsStruct = append(lectureGroupsStruct, lectureGroupStruct)
	}
	return lectureGroupsStruct, nil
}
