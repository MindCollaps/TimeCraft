package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeSlot struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	LecturerID      primitive.ObjectID `json:"lecturerID" bson:"lecturerID"`
	LectureID       primitive.ObjectID `json:"lectureID" bson:"lectureID"`
	TimeStart       primitive.DateTime `json:"timeStart" bson:"timeStart"`
	TimeEnd         primitive.DateTime `json:"timeEnd" bson:"timeEnd"`
	IsOnline        bool               `json:"isOnline" bson:"isOnline"`
	IsReExamination bool               `json:"isReExamination" bson:"isReExamination"`
	IsExam          bool               `json:"isExam" bson:"isExam"`
	IsCancelled     bool               `json:"isCancelled" bson:"isCancelled"`
	WasMoved        bool               `json:"wasMoved" bson:"wasMoved"`
	IsEvent         bool               `json:"isEvent" bson:"isEvent"`
	RoomConfigID    primitive.ObjectID `json:"roomConfigID" bson:"roomConfigID"`
	LastUpdated     primitive.DateTime `json:"lastUpdated" bson:"lastUpdated"`
}
