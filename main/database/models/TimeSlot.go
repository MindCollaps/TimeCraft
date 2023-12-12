package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeSlot struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	LecturerId      primitive.ObjectID `json:"lecturerId" bson:"lecturerId"`
	LectureId       primitive.ObjectID `json:"lectureId" bson:"lectureId"`
	TimeStart       primitive.DateTime `json:"timeStart" bson:"timeStart"`
	TimeEnd         primitive.DateTime `json:"timeEnd" bson:"timeEnd"`
	IsOnline        bool               `json:"isOnline" bson:"isOnline"`
	IsReExamination bool               `json:"isReExamination" bson:"isReExamination"`
	IsExam          bool               `json:"isExam" bson:"isExam"`
	IsCancelled     bool               `json:"isCancelled" bson:"isCancelled"`
	WasMoved        bool               `json:"wasMoved" bson:"wasMoved"`
	IsEvent         bool               `json:"isEvent" bson:"isEvent"`
	RoomConfigId    primitive.ObjectID `json:"roomConfigId" bson:"roomConfigId"`
	LastUpdated     primitive.DateTime `json:"lastUpdated" bson:"lastUpdated"`
}

type TimeSlotStruct struct {
	ID              string     `json:"id" bson:"_id"`
	Name            string     `json:"name" bson:"name"`
	LectureId       Lecture    `json:"lectureId" bson:"lectureId"`
	TimeStart       string     `json:"timeStart" bson:"timeStart"`
	TimeEnd         string     `json:"timeEnd" bson:"timeEnd"`
	IsOnline        bool       `json:"isOnline" bson:"isOnline"`
	IsReExamination bool       `json:"isReExamination" bson:"isReExamination"`
	IsExam          bool       `json:"isExam" bson:"isExam"`
	IsCancelled     bool       `json:"isCancelled" bson:"isCancelled"`
	WasMoved        bool       `json:"wasMoved" bson:"wasMoved"`
	IsEvent         bool       `json:"isEvent" bson:"isEvent"`
	RoomConfigId    RoomConfig `json:"roomConfigId" bson:"roomConfigId"`
	LastUpdated     string     `json:"lastUpdated" bson:"lastUpdated"`
}
