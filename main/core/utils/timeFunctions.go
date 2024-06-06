package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"src/main/env"
	"time"
)

func InitTimeZone() {
	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		log.Println("The 'TIMEZONE' environmental variable is not set. Defaulting to 'Europe/Berlin'.")
		timezone = "Europe/Berlin"
	}
	loc, _ := time.LoadLocation("Europe/Berlin")
	env.TimeZone = loc
}

func ConvertToDateTime(layout string, input string) primitive.DateTime {
	//set timezone to local
	timezone := env.TimeZone
	parsedTime, err := time.ParseInLocation(layout, input, timezone)
	if err != nil {
		log.Println("Error parsing time:", err)
	}
	return primitive.NewDateTimeFromTime(parsedTime)
}

func ConvertToLocalTimeFormat(layout string, input primitive.DateTime) string {
	//set timezone to local
	timezone := env.TimeZone
	return input.Time().In(timezone).Format(layout)
}

func ConvertToLocalTimeObject(input primitive.DateTime) time.Time {
	//set timezone to local
	timezone := env.TimeZone
	return input.Time().In(timezone)
}
