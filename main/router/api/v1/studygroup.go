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
			fmt.Println("Study Group already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// other db query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
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

	cg.GET("/{id}", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Please give correct ID"})
			return
		} else {
			//DB query
		}
	})

	cg.PATCH("/{id}", func(c *gin.Context) {

		//Update
	})

	cg.DELETE("/{id}", func(c *gin.Context) {

		//Delete
	})
}
