package v1

import (
	"bytes"
	"errors"
	"github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"src/main/core/utils"
	"src/main/database"
	"src/main/database/models"
	"time"
)

// /api/v1/ical/...
func icalHandler(cg *gin.RouterGroup) {
	cg.POST("/tbl/:id", func(c *gin.Context) {
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

				start := utils.ConvertToLocalTimeObject(timeslot.TimeStart)
				end := utils.ConvertToLocalTimeObject(timeslot.TimeEnd)
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
			TimeTableId: objectID,
			LastUpdated: utils.ConvertToDateTime(time.DateTime, time.Now().Format(time.DateTime)),
		}
		err = SaveIcalEntry(c, icalEntry)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Failed to save iCal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "iCal generated and saved successfully", "id": icalEntry.ID.Hex()})
	})

	cg.GET("/", func(c *gin.Context) {
		type ListItem struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			TimeTableId string `json:"timeTableId"`
		}

		cursor, err := database.MongoDB.Collection("IcalEntry").Find(c, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Failed to load iCals"})
			log.Println(err)
			return
		}
		defer cursor.Close(c)

		var icalEntries []models.IcalEntry
		if err = cursor.All(c, &icalEntries); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Failed to load iCals"})
			log.Println(err)
			return
		}

		var listItems []ListItem
		for _, icalEntry := range icalEntries {
			listItems = append(listItems, ListItem{
				ID:          icalEntry.ID.Hex(),
				Name:        icalEntry.Name,
				TimeTableId: icalEntry.TimeTableId.Hex(),
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": listItems})
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

	cg.GET("/tbl/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid TimeTable ID"})
			log.Println(err)
			return
		}

		var icalEntry models.IcalEntry
		err = database.MongoDB.Collection("IcalEntry").FindOne(c, bson.M{"timeTableID": id}).Decode(&icalEntry)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "iCal not found"})
			log.Println(err)
			return
		}

		c.Data(http.StatusOK, "text/calendar", []byte(icalEntry.Text))
	})

	cg.PATCH("/:id", func(c *gin.Context) {
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

		var timetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": icalEntry.TimeTableId}).Decode(&timetable)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			log.Println(err)
			return

		}

		var requestBody struct {
			Name        string             `json:"name"`
			Text        string             `json:"text"`
			TimeTableId primitive.ObjectID `json:"timeTableId"`
			LastUpdated primitive.DateTime `json:"lastUpdated"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
			return
		}

		update := bson.M{}
		if requestBody.Name != "" {
			update["name"] = requestBody.Name
		}

		if requestBody.Text != "" {
			// check if the text is a valid ics stream
			reader := bytes.NewReader([]byte(requestBody.Text))
			_, err := ics.ParseCalendar(reader)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid iCal text", "error": "The provided text can't be parsed as ical"})
				log.Println(err)
				return
			}
			update["text"] = requestBody.Text
		}

		if requestBody.TimeTableId != primitive.NilObjectID {
			if err := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": requestBody.TimeTableId}).Decode(&timetable); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
				log.Println(err)
				return
			}
			update["timeTableId"] = requestBody.TimeTableId
		}

		if requestBody.LastUpdated != primitive.DateTime(0) {
			providedTime := requestBody.LastUpdated.Time()
			if providedTime.After(time.Now()) || providedTime.Equal(time.Now()) {
				update["lastUpdated"] = requestBody.LastUpdated
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid lastUpdated"})
				return
			}
		}

		result, err := database.MongoDB.Collection("IcalEntry").UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Warning: No data provided to update the IcalEntry")
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "iCal updated successfully"})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		result, err := database.MongoDB.Collection("IcalEntry").DeleteOne(c, bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "iCal not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "iCal deleted successfully"})
	})

}

func SaveIcalEntry(c *gin.Context, icalEntry models.IcalEntry) error {
	_, err := database.MongoDB.Collection("IcalEntry").InsertOne(c, icalEntry)
	if err != nil {
		return err
	}
	return nil
}
