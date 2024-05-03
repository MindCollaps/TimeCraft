package v1

import (
	"fmt"
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

func sgrpHandler(cg *gin.RouterGroup) {
	//    /api/v1/sgrp/...
	cg.POST("/sgrp", func(c *gin.Context) {
		var requestBody struct {
			Id                 primitive.ObjectID `json:"id" binding:"required"`
			Name               string             `json:"name" binding:"required"`
			StudentGroupIds    []string           `json:"studentGroupIds" binding:"required"`
			TimeTableId        primitive.ObjectID `json:"timeTableId" binding:"required"`
			SpecialisationsIds []string           `json:"specialisationsIds" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name
		studentGroupIds := requestBody.StudentGroupIds
		specialisationsIds := requestBody.SpecialisationsIds

		//to primitive.objectid
		studentGroupIdsPrimitive := make([]primitive.ObjectID, len(studentGroupIds))
		specialisationsIdsPrimitive := make([]primitive.ObjectID, len(specialisationsIds))

		//put studentGroupIds to studentGroupIdsPrimitive
		for i := 0; i < len(studentGroupIds); i++ {
			studentGroupIdsPrimitive[i], _ = primitive.ObjectIDFromHex(studentGroupIds[i])
		}

		//put specialisationsIds to specialisationsIdsPrimitive
		for i := 0; i < len(specialisationsIds); i++ {
			specialisationsIdsPrimitive[i], _ = primitive.ObjectIDFromHex(specialisationsIds[i])
		}

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
			StudentGroupIds:    studentGroupIdsPrimitive,
			TimeTableId:        primitive.NewObjectID(),
			SpecialisationsIds: specialisationsIdsPrimitive,
		}

		database.MongoDB.Collection("SemesterGroup").InsertOne(c, newSemesterGroup, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created Semester Group"})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var semestergroup models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&semestergroup)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "SemesterGroup not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.JSON(http.StatusOK, semestergroup)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("SemesterGroup").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "SemesterGroup not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "SemesterGroup deleted"})
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var existingsgrp models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&existingsgrp)

		var requestBody struct {
			Name               string               `json:"name"`
			StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds"`
			TimeTableId        *primitive.ObjectID  `json:"timeTableId"`
			SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{}
		if requestBody.Name == "" {
			update["name"] = existingsgrp.Name
		} else {
			update["name"] = requestBody.Name
		}
		if requestBody.StudentGroupIds == nil {
			update["studentGroupIds"] = existingsgrp.StudentGroupIds
		} else {
			update["studentGroupIds"] = requestBody.StudentGroupIds
		}
		if requestBody.TimeTableId == nil {
			update["timeTableId"] = existingsgrp.TimeTableId
		} else {
			update["timeTableId"] = requestBody.TimeTableId
		}
		if requestBody.SpecialisationsIds == nil {
			update["specialisationsIds"] = existingsgrp.SpecialisationsIds
		} else {
			update["specialisationsIds"] = requestBody.SpecialisationsIds
		}
		result, err := database.MongoDB.Collection("SemesterGroup").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "SemesterGroup not found"})
			return
		}
		var updatedSemesterGroup models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedSemesterGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, updatedSemesterGroup)
	})
}
