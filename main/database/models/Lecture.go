package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type Lecture struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	ContactHours int                `json:"contactHours" bson:"contactHours"`
}

type LectureStruct struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	ContactHours int    `json:"contactHours" bson:"contactHours"`
}

// Function to convert Lecture to LectureStruct
func LectureToStruct(lecture Lecture) LectureStruct {
	lectureStruct := LectureStruct{
		ID:           lecture.ID.Hex(),
		Name:         lecture.Name,
		ContactHours: lecture.ContactHours,
	}
	return lectureStruct
}

// Function to load multiple lectures and convert them to LectureStruct
func LoadLectures(c *gin.Context, lectureIDs []primitive.ObjectID) ([]LectureStruct, error) {
	var lecturesStruct []LectureStruct
	for _, lectureID := range lectureIDs {
		var lecture Lecture
		err := database.MongoDB.Collection("Lecture").FindOne(c, bson.M{
			"_id": lectureID,
		}).Decode(&lecture)

		if err != nil {
			continue
		}

		lectureStruct := LectureToStruct(lecture)
		lecturesStruct = append(lecturesStruct, lectureStruct)
	}
	return lecturesStruct, nil
}
