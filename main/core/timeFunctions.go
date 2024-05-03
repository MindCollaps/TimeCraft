package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func ConvertToDateTime(layout string, input string) primitive.DateTime {
	//set timezone to local
	loc, _ := time.LoadLocation("Europe/Berlin")
	parsedTime, err := time.ParseInLocation(layout, input, loc)
	if err != nil {
		log.Println("Error parsing time:", err)
	}
	return primitive.DateTime(primitive.NewDateTimeFromTime(parsedTime))
}

func ConvertToLocalTime(layout string, input primitive.DateTime) string {
	//set timezone to local
	loc, _ := time.LoadLocation("Europe/Berlin")
	return input.Time().In(loc).Format(layout)
}
