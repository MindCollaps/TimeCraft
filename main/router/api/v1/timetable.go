package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"src/main/database"
	"src/main/database/models"
)

func tblHandler(cg *gin.RouterGroup) {
	//    /api/v1/tbl/...
	cg.POST("/create", func(c *gin.Context) {
		//check body for username and password
		var requestBody struct {
			Id   primitive.ObjectID   `json:"id" binding:"required"`
			Name string               `json:"name" binding:"required"`
			Days []primitive.ObjectID `json:"days" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name
		days := requestBody.Days

		var existingTbl models.TimeTable
		err := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"Name": name}).Decode(&existingTbl)

		if err == nil {
			// Table with the same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "TimeTable already exists"})
			fmt.Println("TimeTable already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
			return
		}

		newTable := models.TimeTable{
			ID:   primitive.NewObjectID(),
			Name: name,
			Days: days,
		}

		database.MongoDB.Collection("TimeTable").InsertOne(c, newTable, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created Timetable"})

	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var timetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&timetable)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "TimeTable not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.JSON(http.StatusOK, timetable)
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var existingtbl models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&existingtbl)

		var requestBody struct {
			Name string               `json:"name"`
			Days []primitive.ObjectID `json:"days"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{}
		if requestBody.Name == "" {
			update["name"] = existingtbl.Name
		} else {
			update["name"] = requestBody.Name
		}
		if requestBody.Days == nil {
			update["days"] = existingtbl.Days
		} else {
			update["days"] = requestBody.Days
		}
		result, err := database.MongoDB.Collection("TimeTable").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "TimeTable not found"})
			return
		}
		var updatedTimetable models.TimeTable
		err = database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedTimetable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, updatedTimetable)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("TimeTable").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			fmt.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "TimeTable not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "TimeTable deleted"})
	})
}
