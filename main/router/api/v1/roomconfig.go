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

// /api/v1/rmc/...
func rmcHandler(cg *gin.RouterGroup) {
	cg.POST("/", func(c *gin.Context) {

		var requestBody struct {
			Name           string               `json:"name" binding:"required"`
			RoomNr         string               `json:"roomNr" binding:"required"`
			Capacity       int                  `json:"capacity" binding:"required"`
			ExamCapacity   int                  `json:"examCapacity" binding:"required"`
			Blocks         []primitive.ObjectID `json:"blocks" binding:"required"`
			Specialisation []primitive.ObjectID `json:"specialisation" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		name := requestBody.Name
		roomNr := requestBody.RoomNr
		capacity := requestBody.Capacity
		examCapacity := requestBody.ExamCapacity
		blocks := requestBody.Blocks
		specialisation := requestBody.Specialisation

		var existingRmc models.RoomConfig
		err := database.MongoDB.Collection("RoomConfig").FindOne(c, bson.M{"Name": name}).Decode(&existingRmc)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "RoomConfig already exists"})
			log.Println("RoomConfig already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		newRoomConfig := models.RoomConfig{
			ID:             primitive.NewObjectID(),
			Name:           name,
			RoomNr:         roomNr,
			Capacity:       capacity,
			ExamCapacity:   examCapacity,
			Blocks:         blocks,
			Specialisation: specialisation,
		}

		result, err := database.MongoDB.Collection("RoomConfig").InsertOne(c, newRoomConfig, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created RoomConfig", "id": result.InsertedID})

	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		var roomconfig models.RoomConfig
		err = database.MongoDB.Collection("RoomConfig").FindOne(c, bson.M{"_id": objectID}).Decode(&roomconfig)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "RoomConfig not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			}
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, roomconfig)
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}
		var existingrmc models.RoomConfig
		err = database.MongoDB.Collection("RoomConfig").FindOne(c, bson.M{"_id": objectID}).Decode(&existingrmc)

		var requestBody struct {
			Name           string               `json:"name"`
			RoomNr         string               `json:"roomNr"`
			Capacity       int                  `json:"capacity"`
			ExamCapacity   int                  `json:"examCapacity"`
			Blocks         []primitive.ObjectID `json:"blocks"`
			Specialisation []primitive.ObjectID `json:"specialisation"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		update := bson.M{}
		if requestBody.Name != "" {
			update["name"] = requestBody.Name
		}
		if requestBody.RoomNr != "" {
			update["roomNr"] = requestBody.RoomNr
		}
		if requestBody.Capacity != 0 {
			update["capacity"] = requestBody.Capacity
		}
		if requestBody.ExamCapacity != 0 {
			update["examCapacity"] = requestBody.ExamCapacity
		}
		if requestBody.Blocks != nil && !utils.ContainsNilObjectID(requestBody.Blocks) && len(requestBody.Blocks) != 0 {
			update["blocks"] = requestBody.Blocks
		}
		if requestBody.Specialisation != nil && !utils.ContainsNilObjectID(requestBody.Specialisation) && len(requestBody.Specialisation) != 0 {
			update["specialisation"] = requestBody.Specialisation
		}

		result, err := database.MongoDB.Collection("RoomConfig").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Error: RoomConfig not found")
			return
		}

		var updatedRoomConfig models.RoomConfig
		err = database.MongoDB.Collection("RoomConfig").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedRoomConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, updatedRoomConfig)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			log.Println(err)
			return
		}

		result, err := database.MongoDB.Collection("RoomConfig").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "RoomConfig not found"})
			log.Println("Error: RoomConfig not found")
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "RoomConfig deleted"})
	})
}
