package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type LecturersPreferedLecture struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	LectureId  primitive.ObjectID `json:"lectureId" bson:"lectureId"`
	MaxPerWeek int                `json:"maxPerWeek" bson:"maxPerWeek"`
	MaxPerDay  int                `json:"maxPerDay" bson:"maxPerDay"`
	Type       int                `json:"type" bson:"type"`
}

type LecturersPreferedLectureStruct struct {
	ID         string        `json:"id" bson:"_id"`
	LectureId  LectureStruct `json:"lectureId" bson:"lectureId"`
	MaxPerWeek int           `json:"maxPerWeek" bson:"maxPerWeek"`
	MaxPerDay  int           `json:"maxPerDay" bson:"maxPerDay"`
	Type       int           `json:"type" bson:"type"`
}

func LecturersPreferedLectureToStruct(c *gin.Context, preferedLecture LecturersPreferedLecture) (LecturersPreferedLectureStruct, error) {
	preferedLectureStruct := LecturersPreferedLectureStruct{
		ID:         preferedLecture.ID.Hex(),
		MaxPerWeek: preferedLecture.MaxPerWeek,
		MaxPerDay:  preferedLecture.MaxPerDay,
		Type:       preferedLecture.Type,
	}

	// Convert LectureId to LectureStruct
	lecture, err := LoadLectures(c, []primitive.ObjectID{preferedLecture.LectureId})
	if err != nil {
		return LecturersPreferedLectureStruct{}, err
	}

	preferedLectureStruct.LectureId = lecture[0]

	return preferedLectureStruct, nil
}

func LoadLecturersPreferedLectures(c *gin.Context, preferedLectureIDs []primitive.ObjectID) ([]LecturersPreferedLectureStruct, error) {
	var preferedLecturesStruct []LecturersPreferedLectureStruct
	for _, preferedLectureID := range preferedLectureIDs {
		var preferedLecture LecturersPreferedLecture
		err := database.MongoDB.Collection("LecturersPreferedLectures").FindOne(c, bson.M{
			"_id": preferedLectureID,
		}).Decode(&preferedLecture)

		if err != nil {
			continue
		}

		preferedLectureStruct, err := LecturersPreferedLectureToStruct(c, preferedLecture)
		if err != nil {
			continue
		}
		preferedLecturesStruct = append(preferedLecturesStruct, preferedLectureStruct)
	}
	return preferedLecturesStruct, nil
}
