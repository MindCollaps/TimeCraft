package v1

import (
	"errors"
	"github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"src/main/core"
	"src/main/database"
	"src/main/database/models"
	"time"
)

func icalHandler(cg *gin.RouterGroup) {
	// /api/v1/ical/...
	cg.POST("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existingIcal models.IcalEntry
		err = database.MongoDB.Collection("IcalEntry").FindOne(c, bson.M{"name": objectID}).Decode(&existingIcal)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "iCal already exists"})
			log.Println("iCal already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		var timetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&timetable)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			log.Println(err)
			return
		}

		cal := ics.NewCalendar()
		cal.SetMethod(ics.MethodRequest)

		for _, timeTableDayID := range timetable.Days {

			// get the TimeTableDay
			var timeTableDay models.TimeTableDay
			err := database.MongoDB.Collection("TimeTableDay").FindOne(c, bson.M{"_id": timeTableDayID}).Decode(&timeTableDay)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTableDay not found"})
				log.Println(err)
				return
			}

			for _, timeslotID := range timeTableDay.TimeSlotIds {
				event := cal.AddEvent(primitive.NewObjectID().Hex())
				event.SetCreatedTime(time.Now())
				event.SetDtStampTime(time.Now())
				event.SetModifiedAt(time.Now())

				// get the TimeSlot
				var timeslot models.TimeSlot
				err := database.MongoDB.Collection("TimeSlot").FindOne(c, bson.M{"_id": timeslotID}).Decode(&timeslot)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeSlot not found"})
					log.Println(err)
					return
				}

				start := core.ConvertToLocalTimeObject(timeslot.TimeStart)
				end := core.ConvertToLocalTimeObject(timeslot.TimeEnd)
				event.SetStartAt(start)
				event.SetEndAt(end)
				event.SetSummary(timeslot.Name)
			}
		}

		icalString := cal.Serialize()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Failed to generate iCal"})
			return
		}

		icalEntry := models.IcalEntry{
			ID:          primitive.NewObjectID(),
			Name:        timetable.Name,
			Text:        icalString,
			LastUpdated: core.ConvertToDateTime(time.DateTime, time.Now().Format(time.DateTime)),
		}
		err = SaveIcalEntry(c, icalEntry)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Failed to save iCal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "iCal generated and saved successfully", "data": icalEntry.ID.Hex()})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var icalEntry models.IcalEntry
		err = database.MongoDB.Collection("IcalEntry").FindOne(c, bson.M{"_id": id}).Decode(&icalEntry)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "iCal not found"})
			log.Println(err)
			return
		}

		c.Data(http.StatusOK, "text/calendar", []byte(icalEntry.Text))
	})

}

func SaveIcalEntry(c *gin.Context, icalEntry models.IcalEntry) error {
	_, err := database.MongoDB.Collection("IcalEntry").InsertOne(c, icalEntry)
	if err != nil {
		return err
	}
	return nil
}
