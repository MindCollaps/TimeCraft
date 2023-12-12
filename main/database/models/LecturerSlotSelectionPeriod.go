package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID               string     `json:"id" bson:"_id"`
	PeriodBeginn     string     `json:"periodBeginn" bson:"periodBeginn"`
	PeriodEnd        string     `json:"periodEnd" bson:"periodEnd"`
	SlotSelectionIds []TimeSlot `json:"slotSelectionIds" bson:"slotSelectionIds"`
	NotifyEmail      bool       `json:"notifyEmail" bson:"notifyEmail"`
	IsDone           bool       `json:"isDone" bson:"isDone"`
}
