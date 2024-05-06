package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type SplitGroup struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	TimeTableId   primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
	SplitsStudent []primitive.ObjectID `json:"splitsStudent" bson:"splitsStudent"`
	Size          int                  `json:"size" bson:"size"`
}

type SplitGroupStruct struct {
	ID            string               `json:"id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	TimeTableId   TimeTableStruct      `json:"timeTableId" bson:"timeTableId"`
	SplitsStudent []StudentGroupStruct `json:"splitsStudent" bson:"splitsStudent"`
	Size          int                  `json:"size" bson:"size"`
}

func SplitGroupToStruct(c *gin.Context, splitGroup SplitGroup) (SplitGroupStruct, error) {
	splitGroupStruct := SplitGroupStruct{
		ID:          splitGroup.ID.Hex(),
		Name:        splitGroup.Name,
		Size:        splitGroup.Size,
		TimeTableId: TimeTableStruct{},
	}

	// Convert SplitsStudent to []StudentGroupStruct
	splitsStudent, err := LoadStudentGroups(c, splitGroup.SplitsStudent)
	if err != nil {
		return SplitGroupStruct{}, err
	}

	splitGroupStruct.SplitsStudent = splitsStudent

	return splitGroupStruct, nil
}

func LoadSplitGroups(c *gin.Context, splitGroupIDs []primitive.ObjectID) ([]SplitGroupStruct, error) {
	var SplitGroupsStruct []SplitGroupStruct
	for _, splitGroupID := range splitGroupIDs {
		var splitGroup SplitGroup
		err := database.MongoDB.Collection("SplitGroup").FindOne(c, bson.M{
			"_id": splitGroupID,
		}).Decode(&splitGroup)

		if err != nil {
			continue
		}

		splitGroupStruct, err := SplitGroupToStruct(c, splitGroup)
		if err != nil {
			continue
		}
		SplitGroupsStruct = append(SplitGroupsStruct, splitGroupStruct)
	}

	return SplitGroupsStruct, nil
}
