package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type StudentGroup struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id"`
	Name            string               `json:"name" bson:"name"`
	LectureGroupIds []primitive.ObjectID `json:"lectureGroupIds" bson:"lectureGroupIds"`
	TimeTableId     primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
}

type StudentGroupStruct struct {
	ID              string               `json:"id" bson:"_id"`
	Name            string               `json:"name" bson:"name"`
	LectureGroupIds []LectureGroupStruct `json:"lectureGroupIds" bson:"lectureGroupIds"`
	TimeTableId     TimeTableStruct      `json:"timeTableId" bson:"timeTableId"`
}

func StudentGroupToStruct(c *gin.Context, studentGroup StudentGroup) (StudentGroupStruct, error) {
	studentGroupStruct := StudentGroupStruct{
		ID:          studentGroup.ID.Hex(),
		Name:        studentGroup.Name,
		TimeTableId: TimeTableStruct{}, // Initialize with empty TimeTableStruct
	}

	// Convert LectureGroupIds to []LectureGroupStruct
	lectureGroups, err := LoadLectureGroups(c, studentGroup.LectureGroupIds)
	if err != nil {
		return StudentGroupStruct{}, err
	}

	studentGroupStruct.LectureGroupIds = lectureGroups

	// Load TimeTableId
	if !studentGroup.TimeTableId.IsZero() {
		timeTable, err := LoadTimeTables(c, []primitive.ObjectID{studentGroup.TimeTableId})
		if err != nil {
			return StudentGroupStruct{}, err
		}

		if len(timeTable) > 0 {
			studentGroupStruct.TimeTableId = timeTable[0]
		}
	}

	return studentGroupStruct, nil
}

func LoadStudentGroups(c *gin.Context, studentGroupIDs []primitive.ObjectID) ([]StudentGroupStruct, error) {
	var studentGroupsStruct []StudentGroupStruct
	for _, studentGroupID := range studentGroupIDs {
		var studentGroup StudentGroup
		err := database.MongoDB.Collection("StudentGroup").FindOne(c, bson.M{
			"_id": studentGroupID,
		}).Decode(&studentGroup)

		if err != nil {
			continue
		}

		studentGroupStruct, err := StudentGroupToStruct(c, studentGroup)
		if err != nil {
			continue
		}
		studentGroupsStruct = append(studentGroupsStruct, studentGroupStruct)
	}
	return studentGroupsStruct, nil
}
