package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

type SlotSelection struct {
	ID                            primitive.ObjectID   `json:"id" bson:"_id"`
	LectureIds                    []primitive.ObjectID `json:"lectureIds" bson:"lectureIds"`
	BeginnDate                    primitive.DateTime   `json:"beginnDate" bson:"beginnDate"`
	EndDate                       primitive.DateTime   `json:"endDate" bson:"endDate"`
	LecturerId                    primitive.ObjectID   `json:"lecturerId" bson:"lecturerId"`
	PreferedOnline                bool                 `json:"preferedOnline" bson:"preferedOnline"`
	PreferedRoomSpecialisationIds []primitive.ObjectID `json:"preferedRoomSpecialisationIds" bson:"preferedRoomSpecialisationIds"`
	IsBlocked                     bool                 `json:"isBlocked" bson:"isBlocked"`
	Priority                      int                  `json:"priority" bson:"priority"`
}

type SlotSelectionStruct struct {
	ID                            string                     `json:"id" bson:"_id"`
	LectureIds                    []LectureStruct            `json:"lectureIds" bson:"lectureIds"`
	BeginnDate                    string                     `json:"beginnDate" bson:"beginnDate"`
	EndDate                       string                     `json:"endDate" bson:"endDate"`
	LecturerId                    LecturerStruct             `json:"lecturerId" bson:"lecturerId"`
	PreferedOnline                bool                       `json:"preferedOnline" bson:"preferedOnline"`
	PreferedRoomSpecialisationIds []RoomSpecialisationStruct `json:"preferedRoomSpecialisationIds" bson:"preferedRoomSpecialisationIds"`
	IsBlocked                     bool                       `json:"isBlocked" bson:"isBlocked"`
	Priority                      int                        `json:"priority" bson:"priority"`
}

func SlotSelectionToStruct(c *gin.Context, slotSelection SlotSelection) (SlotSelectionStruct, error) {
	slotSelectionStruct := SlotSelectionStruct{
		ID:             slotSelection.ID.Hex(),
		BeginnDate:     slotSelection.BeginnDate.Time().Format(time.DateTime),
		EndDate:        slotSelection.EndDate.Time().Format(time.DateTime),
		PreferedOnline: slotSelection.PreferedOnline,
		IsBlocked:      slotSelection.IsBlocked,
		Priority:       slotSelection.Priority,
	}

	// Convert LectureIds to []LectureStruct
	lectures, err := LoadLectures(c, slotSelection.LectureIds)
	if err != nil {
		return SlotSelectionStruct{}, err
	}

	slotSelectionStruct.LectureIds = lectures

	// Load LecturerId
	if !slotSelection.LecturerId.IsZero() {
		lecturer, err := LoadLecturers(c, []primitive.ObjectID{slotSelection.LecturerId})
		if err != nil {
			return SlotSelectionStruct{}, err
		}
		slotSelectionStruct.LecturerId = lecturer[0]
	}

	// Convert PreferedRoomSpecialisationIds to []RoomSpecialisationStruct
	specialisations, err := LoadRoomSpecialisations(c, slotSelection.PreferedRoomSpecialisationIds)
	if err != nil {
		return SlotSelectionStruct{}, err
	}

	slotSelectionStruct.PreferedRoomSpecialisationIds = specialisations

	return slotSelectionStruct, nil
}

func LoadSlotSelections(c *gin.Context, slotSelectionIDs []primitive.ObjectID) ([]SlotSelectionStruct, error) {
	var slotSelectionsStruct []SlotSelectionStruct
	for _, slotSelectionID := range slotSelectionIDs {
		var slotSelection SlotSelection
		err := database.MongoDB.Collection("SlotSelection").FindOne(c, bson.M{
			"_id": slotSelectionID,
		}).Decode(&slotSelection)

		if err != nil {
			continue
		}

		slotSelectionStruct, err := SlotSelectionToStruct(c, slotSelection)
		if err != nil {
			continue
		}
		slotSelectionsStruct = append(slotSelectionsStruct, slotSelectionStruct)
	}
	return slotSelectionsStruct, nil
}
