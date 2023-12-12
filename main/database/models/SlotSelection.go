package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID                            string               `json:"id" bson:"_id"`
	LectureIds                    []Lecture            `json:"lectureIds" bson:"lectureIds"`
	BeginnDate                    string               `json:"beginnDate" bson:"beginnDate"`
	EndDate                       string               `json:"endDate" bson:"endDate"`
	LecturerId                    Lecturer             `json:"lecturerId" bson:"lecturerId"`
	PreferedOnline                bool                 `json:"preferedOnline" bson:"preferedOnline"`
	PreferedRoomSpecialisationIds []RoomSpecialisation `json:"preferedRoomSpecialisationIds" bson:"preferedRoomSpecialisationIds"`
	IsBlocked                     bool                 `json:"isBlocked" bson:"isBlocked"`
	Priority                      int                  `json:"priority" bson:"priority"`
}
