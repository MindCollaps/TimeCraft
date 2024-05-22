package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
	"time"
)

type TimeTable struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Days        []primitive.ObjectID `json:"days" bson:"days"`
	LastUpdated primitive.DateTime   `json:"lastUpdated" bson:"lastUpdated"`
}

type TimeTableStruct struct {
	ID          string               `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Days        []TimeTableDayStruct `json:"days" bson:"days"`
	LastUpdated string               `json:"lastUpdated" bson:"lastUpdated"`
}

func TimeTableToStruct(c *gin.Context, timeTable TimeTable) (TimeTableStruct, error) {
	timeTableStruct := TimeTableStruct{
		ID:          timeTable.ID.Hex(),
		Name:        timeTable.Name,
		Days:        []TimeTableDayStruct{},
		LastUpdated: timeTable.LastUpdated.Time().Format(time.DateTime),
	}

	// Load TimeTableDays
	days, err := LoadTimeTableDays(c, timeTable.Days)
	if err != nil {
		return TimeTableStruct{}, err
	}

	timeTableStruct.Days = days

	return timeTableStruct, nil
}

func LoadTimeTables(c *gin.Context, timeTableIDs []primitive.ObjectID) ([]TimeTableStruct, error) {
	var timeTablesStruct []TimeTableStruct
	for _, timeTableID := range timeTableIDs {
		var timeTable TimeTable
		err := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{
			"_id": timeTableID,
		}).Decode(&timeTable)

		if err != nil {
			continue
		}

		timeTableStruct, err := TimeTableToStruct(c, timeTable)
		if err != nil {
			continue
		}
		timeTablesStruct = append(timeTablesStruct, timeTableStruct)
	}
	return timeTablesStruct, nil
}
