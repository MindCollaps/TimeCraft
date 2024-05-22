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
	"src/main/database"
	"src/main/database/models"
)

// /api/v1/tsl/...
func tslHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {

		var requestBody struct {
			Name            string             `json:"name" binding:"required"`
			LecturerId      primitive.ObjectID `json:"lecturerId" binding:"required"`
			LectureId       primitive.ObjectID `json:"lectureId" binding:"required"`
			TimeStart       primitive.DateTime `json:"timeStart" binding:"required"`
			TimeEnd         primitive.DateTime `json:"timeEnd" binding:"required"`
			IsOnline        bool               `json:"isOnline" binding:"required"`
			IsReExamination bool               `json:"isReExamination" binding:"required"`
			IsExam          bool               `json:"isExam" binding:"required"`
			IsCancelled     bool               `json:"isCancelled" binding:"required"`
			WasMoved        bool               `json:"wasMoved" binding:"required"`
			IsEvent         bool               `json:"isEvent" binding:"required"`
			IsHoliday       bool               `json:"isHoliday" binding:"required"`
			RoomConfigId    primitive.ObjectID `json:"roomConfigId" binding:"required"`
			LastUpdated     primitive.DateTime `json:"lastUpdated"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		name := requestBody.Name
		lecturerId := requestBody.LecturerId
		lectureId := requestBody.LectureId
		timeStart := requestBody.TimeStart
		timeEnd := requestBody.TimeEnd
		isOnline := requestBody.IsOnline
		isReExamination := requestBody.IsReExamination
		isExam := requestBody.IsExam
		isCancelled := requestBody.IsCancelled
		wasMoved := requestBody.WasMoved
		isEvent := requestBody.IsEvent
		isHoliday := requestBody.IsHoliday
		roomConfigId := requestBody.RoomConfigId

		var existingTsl models.TimeSlot
		err := database.MongoDB.Collection("TimeSlot").FindOne(c, bson.M{"Name": name}).Decode(&existingTsl)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "TimeSlot already exists"})
			log.Println("TimeSlot already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newTimeSlot := models.TimeSlot{
			ID:              primitive.NewObjectID(),
			Name:            name,
			LecturerId:      lecturerId,
			LectureId:       lectureId,
			TimeStart:       timeStart,
			TimeEnd:         timeEnd,
			IsOnline:        isOnline,
			IsReExamination: isReExamination,
			IsExam:          isExam,
			IsCancelled:     isCancelled,
			WasMoved:        wasMoved,
			IsEvent:         isEvent,
			IsHoliday:       isHoliday,
			RoomConfigId:    roomConfigId,
		}

		result, err := database.MongoDB.Collection("TimeSlot").InsertOne(c, newTimeSlot, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Created TimeSlot", "id": result.InsertedID})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var timeslot models.TimeSlot
		err = database.MongoDB.Collection("TimeSlot").FindOne(c, bson.M{"_id": objectID}).Decode(&timeslot)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeSlot not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, timeslot)
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			return
		}

		var requestBody struct {
			Name            string              `json:"name"`
			LecturerId      *primitive.ObjectID `json:"lecturerId"`
			LectureId       *primitive.ObjectID `json:"lectureId"`
			TimeStart       primitive.DateTime  `json:"timeStart"`
			TimeEnd         primitive.DateTime  `json:"timeEnd"`
			IsOnline        bool                `json:"isOnline"`
			IsReExamination bool                `json:"isReExamination"`
			IsExam          bool                `json:"isExam"`
			IsCancelled     bool                `json:"isCancelled"`
			WasMoved        bool                `json:"wasMoved"`
			IsEvent         bool                `json:"isEvent"`
			IsHoliday       bool                `json:"isHoliday"`
			RoomConfigId    *primitive.ObjectID `json:"roomConfigId"`
			LastUpdated     primitive.DateTime  `json:"lastUpdated"`
		}

		var existingtsl models.TimeSlot
		err = database.MongoDB.Collection("TimeSlot").FindOne(c, bson.M{"_id": objectID}).Decode(&existingtsl)

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
			return
		}

		update := bson.M{}
		if requestBody.Name != "" {
			update["name"] = requestBody.Name
		}
		if requestBody.LecturerId != nil && *requestBody.LecturerId != primitive.NilObjectID {
			update["lecturerId"] = requestBody.LecturerId
		}
		if requestBody.LectureId != nil && *requestBody.LectureId != primitive.NilObjectID {
			update["lectureId"] = requestBody.LectureId
		}
		if requestBody.TimeStart != 0 {
			update["timeStart"] = requestBody.TimeStart
		}
		if requestBody.TimeEnd != 0 {
			update["timeEnd"] = requestBody.TimeEnd
		}
		if requestBody.IsOnline {
			update["isOnline"] = requestBody.IsOnline
		}
		if requestBody.IsReExamination {
			update["isReExamination"] = requestBody.IsReExamination
		}
		if requestBody.IsExam {
			update["isExam"] = requestBody.IsExam
		}
		if requestBody.IsCancelled {
			update["isCancelled"] = requestBody.IsCancelled
		}
		if requestBody.WasMoved {
			update["wasMoved"] = requestBody.WasMoved
		}
		if requestBody.IsEvent {
			update["isEvent"] = requestBody.IsEvent
		}
		if requestBody.IsHoliday {
			update["isHoliday"] = requestBody.IsHoliday
		}
		if requestBody.RoomConfigId != nil && *requestBody.RoomConfigId != primitive.NilObjectID {
			update["roomConfigId"] = requestBody.RoomConfigId
		}

		result, err := database.MongoDB.Collection("TimeSlot").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeSlot not found"})
			log.Println("TimeSlot not found")
			return
		}

		var updatedTimeslot models.TimeSlot
		err = database.MongoDB.Collection("TimeSlot").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedTimeslot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, updatedTimeslot)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		result, err := database.MongoDB.Collection("TimeSlot").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeSlot not found"})
			log.Println("Error: TimeSlot not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "TimeSlot deleted"})
	})
}
