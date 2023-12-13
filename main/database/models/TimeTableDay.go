package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

type TimeTableDay struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Date        primitive.DateTime   `json:"date" bson:"date"`
	TimeSlotIds []primitive.ObjectID `json:"timeSlotIds" bson:"timeSlotIds"`
	LastUpdated primitive.DateTime   `json:"lastUpdated" bson:"lastUpdated"`
}

type TimeTableDayStruct struct {
	ID          string           `json:"id" bson:"_id"`
	Date        string           `json:"date" bson:"date"`
	TimeSlotIds []TimeSlotStruct `json:"timeSlotIds" bson:"timeSlotIds"`
	LastUpdated string           `json:"lastUpdated" bson:"lastUpdated"`
}

func TimeTableDayToStruct(c *gin.Context, timeTableDay TimeTableDay) (TimeTableDayStruct, error) {
	timeTableDayStruct := TimeTableDayStruct{
		ID:          timeTableDay.ID.Hex(),
		Date:        timeTableDay.Date.Time().Format(time.DateTime),
		LastUpdated: timeTableDay.LastUpdated.Time().Format(time.DateTime),
	}

	// Load TimeSlotIds
	timeSlots, err := LoadTimeSlots(c, timeTableDay.TimeSlotIds)
	if err != nil {
		return TimeTableDayStruct{}, err
	}

	timeTableDayStruct.TimeSlotIds = timeSlots

	return timeTableDayStruct, nil
}

func LoadTimeTableDays(c *gin.Context, timeTableDayIDs []primitive.ObjectID) ([]TimeTableDayStruct, error) {
	var timeTableDaysStruct []TimeTableDayStruct
	for _, timeTableDayID := range timeTableDayIDs {
		var timeTableDay TimeTableDay
		err := database.MongoDB.Collection("TimeTableDays").FindOne(c, bson.M{
			"_id": timeTableDayID,
		}).Decode(&timeTableDay)

		if err != nil {
			continue
		}

		timeTableDayStruct, err := TimeTableDayToStruct(c, timeTableDay)
		if err != nil {
			continue
		}
		timeTableDaysStruct = append(timeTableDaysStruct, timeTableDayStruct)
	}
	return timeTableDaysStruct, nil
}
