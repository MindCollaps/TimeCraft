package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

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
	ID              string           `json:"id" bson:"_id"`
	Name            string           `json:"name" bson:"name"`
	LectureId       LectureStruct    `json:"lectureId" bson:"lectureId"`
	TimeStart       string           `json:"timeStart" bson:"timeStart"`
	TimeEnd         string           `json:"timeEnd" bson:"timeEnd"`
	IsOnline        bool             `json:"isOnline" bson:"isOnline"`
	IsReExamination bool             `json:"isReExamination" bson:"isReExamination"`
	IsExam          bool             `json:"isExam" bson:"isExam"`
	IsCancelled     bool             `json:"isCancelled" bson:"isCancelled"`
	WasMoved        bool             `json:"wasMoved" bson:"wasMoved"`
	IsEvent         bool             `json:"isEvent" bson:"isEvent"`
	RoomConfigId    RoomConfigStruct `json:"roomConfigId" bson:"roomConfigId"`
	LastUpdated     string           `json:"lastUpdated" bson:"lastUpdated"`
}

func TimeSlotToStruct(c *gin.Context, timeSlot TimeSlot) (TimeSlotStruct, error) {
	timeSlotStruct := TimeSlotStruct{
		ID:              timeSlot.ID.Hex(),
		Name:            timeSlot.Name,
		TimeStart:       timeSlot.TimeStart.Time().Format(time.DateTime),
		TimeEnd:         timeSlot.TimeEnd.Time().Format(time.DateTime),
		IsOnline:        timeSlot.IsOnline,
		IsReExamination: timeSlot.IsReExamination,
		IsExam:          timeSlot.IsExam,
		IsCancelled:     timeSlot.IsCancelled,
		WasMoved:        timeSlot.WasMoved,
		IsEvent:         timeSlot.IsEvent,
		LastUpdated:     timeSlot.LastUpdated.Time().Format(time.DateTime),
	}

	// Load LectureId
	if !timeSlot.LectureId.IsZero() {
		lecture, err := LoadLectures(c, []primitive.ObjectID{timeSlot.LectureId})
		if err != nil {
			return TimeSlotStruct{}, err
		}
		timeSlotStruct.LectureId = lecture[0]
	}

	// Load RoomConfigId
	if !timeSlot.RoomConfigId.IsZero() {
		roomConfig, err := LoadRoomConfigs(c, []primitive.ObjectID{timeSlot.RoomConfigId})
		if err != nil {
			return TimeSlotStruct{}, err
		}
		timeSlotStruct.RoomConfigId = roomConfig[0]
	}

	return timeSlotStruct, nil
}

func LoadTimeSlots(c *gin.Context, timeSlotIDs []primitive.ObjectID) ([]TimeSlotStruct, error) {
	var timeSlotsStruct []TimeSlotStruct
	for _, timeSlotID := range timeSlotIDs {
		var timeSlot TimeSlot
		err := database.MongoDB.Collection("TimeSlots").FindOne(c, bson.M{
			"_id": timeSlotID,
		}).Decode(&timeSlot)

		if err != nil {
			continue
		}

		timeSlotStruct, err := TimeSlotToStruct(c, timeSlot)
		if err != nil {
			continue
		}
		timeSlotsStruct = append(timeSlotsStruct, timeSlotStruct)
	}
	return timeSlotsStruct, nil
}
