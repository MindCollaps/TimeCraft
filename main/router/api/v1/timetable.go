package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"src/main/core"
	"src/main/database"
	"src/main/database/models"
	"time"
)

// /api/v1/tbl/...
func tblHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {
		var requestBody struct {
			Name string               `json:"name" binding:"required"`
			Days []primitive.ObjectID `json:"days" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
			return
		}

		name := requestBody.Name
		days := requestBody.Days

		var existingTbl models.TimeTable
		err := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"name": name}).Decode(&existingTbl)

		if err == nil {
			// Table with the same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "TimeTable already exists"})
			log.Println("TimeTable already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newTable := models.TimeTable{
			ID:   primitive.NewObjectID(),
			Name: name,
			Days: days,
		}

		result, err := database.MongoDB.Collection("TimeTable").InsertOne(c, newTable, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Created Timetable", "id": result.InsertedID})
	})

	cg.GET("/", func(c *gin.Context) {
		type TimeTableResponse struct {
			ID   primitive.ObjectID `json:"id" bson:"_id"`
			Name string             `json:"name" bson:"name"`
		}

		var timetables []TimeTableResponse
		opts := options.Find().SetProjection(bson.M{"name": 1, "_id": 1})
		cursor, err := database.MongoDB.Collection("TimeTable").Find(c, bson.M{}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if err = cursor.All(c, &timetables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Fetched timetables", "data": timetables})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}
		var timetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&timetable)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Fetched timetable", "data": timetable})
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existingtbl models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&existingtbl)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			log.Println(err)
			return
		}

		var requestBody struct {
			Name string               `json:"name"`
			Days []primitive.ObjectID `json:"days"`
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

		if requestBody.Days != nil && !core.ContainsNilObjectID(requestBody.Days) && len(requestBody.Days) != 0 {
			update["days"] = requestBody.Days
		}

		if len(update) > 0 {
			lastUpdated := core.ConvertToDateTime(time.DateTime, time.Now().Format(time.DateTime))
			update["lastUpdated"] = lastUpdated
		}

		result, err := database.MongoDB.Collection("TimeTable").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Warning: No data provided to update the TimeTable")
			return
		}

		var updatedTimetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedTimetable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "TimeTable updated successfully", "data": updatedTimetable})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		// Delete all TimeTableDays and TimeSlots associated with the TimeTable
		timeTableDaysIDs := getAllTimeTableDays(objectID)
		for _, dayID := range timeTableDaysIDs {
			for _, timeSlotID := range getAllTimeSlots(dayID) {
				deleteTimeSlot(timeSlotID)
			}
			deleteTimeTableDay(dayID)
		}

		result, err := database.MongoDB.Collection("TimeTable").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			log.Println("Error: TimeTable not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "TimeTable deleted"})
	})
}
