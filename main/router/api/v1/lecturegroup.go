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
	"src/main/core/utils"
	"src/main/database"
	"src/main/database/models"
)

// /api/v1/lgrp/...
func lgrpHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {
		//check body for name and timeTableId
		var requestBody struct {
			Name        string             `json:"name" binding:"required"`
			TimeTableId primitive.ObjectID `json:"timeTableId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
			return
		}

		name := requestBody.Name
		timeTableId := requestBody.TimeTableId

		var existinglgrp models.LectureGroup
		err := database.MongoDB.Collection("LectureGroup").FindOne(c, bson.M{"Name": name}).Decode(&existinglgrp)

		if err == nil {
			// LectureGroup with the same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "LectureGroup already exists"})
			log.Println("LectureGroup already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newLectureGroup := models.LectureGroup{
			ID:          primitive.NewObjectID(),
			Name:        name,
			TimeTableId: timeTableId,
		}

		result, err := database.MongoDB.Collection("LectureGroup").InsertOne(c, newLectureGroup, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Created LectureGroup", "id": result.InsertedID})

	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var lecturegroup models.LectureGroup
		err = database.MongoDB.Collection("LectureGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&lecturegroup)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "LectureGroup not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "LectureGroup found", "data": lecturegroup})
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existinglgrp models.LectureGroup
		err = database.MongoDB.Collection("LectureGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&existinglgrp)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "LectureGroup not found"})
			log.Println(err)
			return
		}

		var requestBody struct {
			Name        string               `json:"name"`
			TimeTableId []primitive.ObjectID `json:"timeTableId"`
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
		if requestBody.TimeTableId != nil && len(requestBody.TimeTableId) > 0 && !utils.ContainsNilObjectID(requestBody.TimeTableId) {
			update["timeTableId"] = requestBody.TimeTableId
		}

		result, err := database.MongoDB.Collection("LectureGroup").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Warning: No data provided to update the LectureGroup")
			return
		}

		var updatedLectureGroup models.LectureGroup
		err = database.MongoDB.Collection("LectureGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedLectureGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "LectureGroup updated", "data": updatedLectureGroup})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}
		result, err := database.MongoDB.Collection("LectureGroup").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "LectureGroup not found"})
			log.Println("Error: LectureGroup not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "LectureGroup deleted"})
	})
}
