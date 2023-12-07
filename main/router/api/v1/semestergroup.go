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

func sgrpHandler(cg *gin.RouterGroup) {
	//    /api/v1/sgrp/...
	cg.POST("/sgrp", func(c *gin.Context) {
		var requestBody struct {
			Id                 primitive.ObjectID `json:"id" binding:"required"`
			Name               string             `json:"name" binding:"required"`
			StudentGroupIds    string             `json:"studentGroupIds" binding:"required"`
			TimeTableId        primitive.ObjectID `json:"timeTableId" binding:"required"`
			SpecialisationsIds string             `json:"specialisationsIds" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name
		studentGroupIds := requestBody.StudentGroupIds
		specialisationsIds := requestBody.SpecialisationsIds

		var existingSgrp models.SemesterGroup
		err := database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"Name": name}).Decode(&existingSgrp)

		if err == nil {
			// Semester Group with same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "Semester Group already exists"})
			fmt.Println("Semester Group already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// other db query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
			return
		}

		newSemesterGroup := models.SemesterGroup{
			ID:                 primitive.NewObjectID(),
			Name:               name,
			StudentGroupIds:    studentGroupIds,
			TimeTableID:        primitive.NewObjectID(),
			SpecialisationsIds: specialisationsIds,
		}

		database.MongoDB.Collection("SemesterGroup").InsertOne(c, newSemesterGroup, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"status": 200, "message": "Created Semester Group"})
	})
}
