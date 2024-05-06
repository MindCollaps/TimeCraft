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

// /api/v1/sgrp/...
func sgrpHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {
		var requestBody struct {
			Name               string               `json:"name" binding:"required"`
			StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds" binding:"required"`
			TimeTableId        primitive.ObjectID   `json:"timeTableId" binding:"required"`
			SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		name := requestBody.Name
		studentGroupIds := requestBody.StudentGroupIds
		specialisationsIds := requestBody.SpecialisationsIds

		var existingSgrp models.SemesterGroup
		err := database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"Name": name}).Decode(&existingSgrp)

		if err == nil {
			// Semester Group with same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "Semester Group already exists"})
			log.Println("Semester Group already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// other db query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newSemesterGroup := models.SemesterGroup{
			ID:                 primitive.NewObjectID(),
			Name:               name,
			StudentGroupIds:    studentGroupIds,
			TimeTableId:        primitive.NewObjectID(),
			SpecialisationsIds: specialisationsIds,
		}

		_, err = database.MongoDB.Collection("SemesterGroup").InsertOne(c, newSemesterGroup, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Created Semester Group"})
	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}
		var semestergroup models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&semestergroup)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SemesterGroup not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, semestergroup)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}
		result, err := database.MongoDB.Collection("SemesterGroup").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SemesterGroup not found"})
			log.Println("Error: SemesterGroup not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "SemesterGroup deleted"})
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var existingsgrp models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&existingsgrp)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SemesterGroup not found"})
			log.Println(err)
			return
		}

		var requestBody struct {
			Name               string               `json:"name"`
			StudentGroupIds    []primitive.ObjectID `json:"studentGroupIds"`
			TimeTableId        *primitive.ObjectID  `json:"timeTableId"`
			SpecialisationsIds []primitive.ObjectID `json:"specialisationsIds"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid body"})
			log.Println(err)
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
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "SemesterGroup not found"})
			log.Println(err)
			return
		}
		var updatedSemesterGroup models.SemesterGroup
		err = database.MongoDB.Collection("SemesterGroup").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedSemesterGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, updatedSemesterGroup)
	})
}
