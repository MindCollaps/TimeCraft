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
)

// /api/v1/splitgrp/...
func splitgrpHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {
		var requestBody struct {
			Name          string               `json:"name" binding:"required"`
			TimeTableId   primitive.ObjectID   `json:"timeTableId" binding:"required"`
			SplitsStudent []primitive.ObjectID `json:"splitsStudent" binding:"required"`
			Size          int                  `json:"size" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
			return
		}

		// Check if the SplitGroup already exists
		var existingSplitGroup models.SplitGroup
		err := database.MongoDB.Collection("SplitGroup").FindOne(c, bson.M{"name": requestBody.Name}).Decode(&existingSplitGroup)
		if err == nil {
			// SplitGroup with the same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "SplitGroup already exists"})
			log.Println("SplitGroup already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newSplitGroup := models.SplitGroup{
			ID:            primitive.NewObjectID(),
			Name:          requestBody.Name,
			TimeTableId:   requestBody.TimeTableId,
			SplitsStudent: requestBody.SplitsStudent,
			Size:          requestBody.Size,
		}

		_, err = database.MongoDB.Collection("SplitGroup").InsertOne(c, newSplitGroup, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Successfully created the SplitGroup"})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existingSplitGroup models.SplitGroup
		err = database.MongoDB.Collection("SplitGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&existingSplitGroup)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SplitGroup not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, existingSplitGroup)
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existingSplitGroup models.SplitGroup
		err = database.MongoDB.Collection("SplitGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&existingSplitGroup)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "TimeTable not found"})
			log.Println(err)
			return
		}

		var requestBody struct {
			Name          string               `json:"name"`
			TimeTableId   primitive.ObjectID   `json:"timeTableId"`
			SplitsStudent []primitive.ObjectID `json:"splitsStudent"`
			Size          int                  `json:"size"`
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

		if requestBody.TimeTableId != primitive.NilObjectID {
			update["timeTableId"] = requestBody.TimeTableId
		}

		if requestBody.SplitsStudent != nil && len(requestBody.SplitsStudent) != 0 && !core.ContainsNilObjectID(requestBody.SplitsStudent) {
			update["splitsStudent"] = requestBody.SplitsStudent
		}

		if requestBody.Size >= 1 {
			update["size"] = requestBody.Size
		}

		result, err := database.MongoDB.Collection("SplitGroup").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Warning: No data provided to update the SplitGroup")
			return
		}

		var updatedSplitGroup models.SplitGroup
		err = database.MongoDB.Collection("SplitGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedSplitGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, updatedSplitGroup)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		result, err := database.MongoDB.Collection("SplitGroup").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SplitGroup not found"})
			log.Println("Error: SplitGroup not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "SplitGroup deleted"})
	})
}
