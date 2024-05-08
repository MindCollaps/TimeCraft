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
}
