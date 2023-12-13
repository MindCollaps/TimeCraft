package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

type LecturerSlotSelection struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	PeriodBeginn     primitive.DateTime   `json:"periodBeginn" bson:"periodBeginn"`
	PeriodEnd        primitive.DateTime   `json:"periodEnd" bson:"periodEnd"`
	SlotSelectionIds []primitive.ObjectID `json:"slotSelectionIds" bson:"slotSelectionIds"`
	NotifyEmail      bool                 `json:"notifyEmail" bson:"notifyEmail"`
	IsDone           bool                 `json:"isDone" bson:"isDone"`
}

type LecturerSlotSelectionStruct struct {
	ID               string           `json:"id" bson:"_id"`
	PeriodBeginn     string           `json:"periodBeginn" bson:"periodBeginn"`
	PeriodEnd        string           `json:"periodEnd" bson:"periodEnd"`
	SlotSelectionIds []TimeSlotStruct `json:"slotSelectionIds" bson:"slotSelectionIds"`
	NotifyEmail      bool             `json:"notifyEmail" bson:"notifyEmail"`
	IsDone           bool             `json:"isDone" bson:"isDone"`
}

func LecturerSlotSelectionToStruct(c *gin.Context, lecturerSlotSelection LecturerSlotSelection) (LecturerSlotSelectionStruct, error) {
	lecturerSlotSelectionStruct := LecturerSlotSelectionStruct{
		ID:           lecturerSlotSelection.ID.Hex(),
		PeriodBeginn: lecturerSlotSelection.PeriodBeginn.Time().Format(time.DateTime),
		PeriodEnd:    lecturerSlotSelection.PeriodEnd.Time().Format(time.DateTime),
		NotifyEmail:  lecturerSlotSelection.NotifyEmail,
		IsDone:       lecturerSlotSelection.IsDone,
	}

	// Convert SlotSelectionIds to []TimeSlotStruct
	slotSelection, err := LoadTimeSlots(c, lecturerSlotSelection.SlotSelectionIds)
	if err != nil {
		return LecturerSlotSelectionStruct{}, err
	}

	lecturerSlotSelectionStruct.SlotSelectionIds = slotSelection

	return lecturerSlotSelectionStruct, nil
}

func LoadLecturerSlotSelections(c *gin.Context, lecturerSlotSelectionIDs []primitive.ObjectID) ([]LecturerSlotSelectionStruct, error) {
	var lecturerSlotSelectionsStruct []LecturerSlotSelectionStruct
	for _, lecturerSlotSelectionID := range lecturerSlotSelectionIDs {
		var lecturerSlotSelection LecturerSlotSelection
		err := database.MongoDB.Collection("LecturerSlotSelection").FindOne(c, bson.M{
			"_id": lecturerSlotSelectionID,
		}).Decode(&lecturerSlotSelection)

		if err != nil {
			continue
		}

		lecturerSlotSelectionStruct, err := LecturerSlotSelectionToStruct(c, lecturerSlotSelection)
		if err != nil {
			continue
		}
		lecturerSlotSelectionsStruct = append(lecturerSlotSelectionsStruct, lecturerSlotSelectionStruct)
	}
	return lecturerSlotSelectionsStruct, nil
}
