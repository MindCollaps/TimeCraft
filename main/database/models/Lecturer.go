package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type Lecturer struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	FirstName    string               `json:"firstName" bson:"firstName"`
	SureName     string               `json:"lastName" bson:"lastName"`
	Title        string               `json:"title" bson:"title"`
	CanHold      []primitive.ObjectID `json:"canHold" bson:"canHold"`
	ContactEmail string               `json:"contactEmail" bson:"contactEmail"`
	ContactPhone string               `json:"contactPhone" bson:"contactPhone"`
	MaxWeekHours int                  `json:"maxWeekHours" bson:"maxWeekHours"`
}

type LecturerStruct struct {
	ID           string          `json:"id" bson:"_id"`
	FirstName    string          `json:"firstName" bson:"firstName"`
	SureName     string          `json:"lastName" bson:"lastName"`
	Title        string          `json:"title" bson:"title"`
	CanHold      []LectureStruct `json:"canHold" bson:"canHold"`
	ContactEmail string          `json:"contactEmail" bson:"contactEmail"`
	ContactPhone string          `json:"contactPhone" bson:"contactPhone"`
	MaxWeekHours int             `json:"maxWeekHours" bson:"maxWeekHours"`
}

func LecturerToStruct(c *gin.Context, lecturer Lecturer) (LecturerStruct, error) {
	lecturerStruct := LecturerStruct{
		ID:           lecturer.ID.Hex(),
		FirstName:    lecturer.FirstName,
		SureName:     lecturer.SureName,
		Title:        lecturer.Title,
		ContactEmail: lecturer.ContactEmail,
		ContactPhone: lecturer.ContactPhone,
		MaxWeekHours: lecturer.MaxWeekHours,
	}

	// Convert CanHold to []LectureStruct
	canHold, err := LoadLectures(c, lecturer.CanHold)
	if err != nil {
		return LecturerStruct{}, err
	}

	lecturerStruct.CanHold = canHold

	return lecturerStruct, nil
}

func LoadLecturers(c *gin.Context, lecturerIDs []primitive.ObjectID) ([]LecturerStruct, error) {
	var lecturersStruct []LecturerStruct
	for _, lecturerID := range lecturerIDs {
		var lecturer Lecturer
		err := database.MongoDB.Collection("Lecturer").FindOne(c, bson.M{
			"_id": lecturerID,
		}).Decode(&lecturer)

		if err != nil {
			continue
		}

		lecturerStruct, err := LecturerToStruct(c, lecturer)
		if err != nil {
			continue
		}
		lecturersStruct = append(lecturersStruct, lecturerStruct)
	}
	return lecturersStruct, nil
}
