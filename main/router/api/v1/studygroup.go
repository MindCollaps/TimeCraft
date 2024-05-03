package v1

import (
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

func stygrpHandler(cg *gin.RouterGroup) {
	//    /api/v1/stygrp/...
	cg.POST("/stygrp", func(c *gin.Context) {
		var requestBody struct {
			Id              primitive.ObjectID `json:"id" binding:"required"`
			Name            string             `json:"name" binding:"required"`
			LectureGroupIds []string           `json:"lectureGroupIds" binding:"required"`
			TimeTableId     primitive.ObjectID `json:"timeTableId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name
		lectureGroupIds := requestBody.LectureGroupIds

		//to primitive.objectid
		lectureGroupIdsPrimitive := make([]primitive.ObjectID, len(lectureGroupIds))

		//put lectureGroupIds to lectureGroupIdsPrimitive
		for i := 0; i < len(lectureGroupIds); i++ {
			lectureGroupIdsPrimitive[i], _ = primitive.ObjectIDFromHex(lectureGroupIds[i])
		}

		// is StudentGroup==StudyGroup???
		var existingStygrp models.StudentGroup
		err := database.MongoDB.Collection("StudyGroup").FindOne(c, bson.M{"Name": name}).Decode(&existingStygrp)

		if err == nil {
			// Study Group with same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "Study Group already exists"})
			log.Println("Study Group already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// other db query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			log.Println(err)
			return
		}

		newStudyGroup := models.StudentGroup{
			ID:              primitive.NewObjectID(),
			Name:            name,
			LectureGroupIds: lectureGroupIdsPrimitive,
			TimeTableId:     primitive.NewObjectID(),
		}

		database.MongoDB.Collection("StudyGroup").InsertOne(c, newStudyGroup, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created Semester Group"})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var studentgroup models.StudentGroup
		err = database.MongoDB.Collection("StudentGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&studentgroup)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "StudentGroup not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.JSON(http.StatusOK, studentgroup)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("StudentGroup").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "StudentGroup not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "StudentGroup deleted"})
	})
}
