package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type SemesterGroup struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Name               string               `json:"name" bson:"name"`
	StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableId        primitive.ObjectID   `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds" bson:"specialisationsIds"`
}

type SemesterGroupStruct struct {
	ID                 string               `json:"id" bson:"_id"`
	Name               string               `json:"name" bson:"name"`
	StudentGroupIds    []StudentGroupStruct `json:"studentGroupIds" bson:"studentGroupIds"`
	TimeTableId        TimeTableStruct      `json:"timeTableId" bson:"timeTableId"`
	SpecialisationsIds []StudentGroupStruct `json:"specialisationsIds" bson:"specialisationsIds"`
}

func SemesterGroupToStruct(c *gin.Context, semesterGroup SemesterGroup) (SemesterGroupStruct, error) {
	semesterGroupStruct := SemesterGroupStruct{
		ID:          semesterGroup.ID.Hex(),
		Name:        semesterGroup.Name,
		TimeTableId: TimeTableStruct{}, // Initialize with empty TimeTableStruct
	}

	// Convert StudentGroupIds to []StudentGroupStruct
	studentGroups, err := LoadStudentGroups(c, semesterGroup.StudentGroupIds)
	if err != nil {
		return SemesterGroupStruct{}, err
	}

	semesterGroupStruct.StudentGroupIds = studentGroups

	// Convert SpecialisationsIds to []RoomSpecialisationStruct
	specialisations, err := LoadStudentGroups(c, semesterGroup.SpecialisationsIds)
	if err != nil {
		return SemesterGroupStruct{}, err
	}

	semesterGroupStruct.SpecialisationsIds = specialisations

	// Load TimeTableId
	if !semesterGroup.TimeTableId.IsZero() {
		timeTable, err := LoadTimeTables(c, []primitive.ObjectID{semesterGroup.TimeTableId})
		if err != nil {
			return SemesterGroupStruct{}, err
		}

		if len(timeTable) > 0 {
			semesterGroupStruct.TimeTableId = timeTable[0]
		}
	}

	return semesterGroupStruct, nil
}

func LoadSemesterGroups(c *gin.Context, semesterGroupIDs []primitive.ObjectID) ([]SemesterGroupStruct, error) {
	var semesterGroupsStruct []SemesterGroupStruct
	for _, semesterGroupID := range semesterGroupIDs {
		var semesterGroup SemesterGroup
		err := database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{
			"_id": semesterGroupID,
		}).Decode(&semesterGroup)

		if err != nil {
			continue
		}

		semesterGroupStruct, err := SemesterGroupToStruct(c, semesterGroup)
		if err != nil {
			continue
		}
		semesterGroupsStruct = append(semesterGroupsStruct, semesterGroupStruct)
	}
	return semesterGroupsStruct, nil
}
