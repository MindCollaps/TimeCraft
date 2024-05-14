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

func rmsHandler(cg *gin.RouterGroup) {
	//    /api/v1/rms/...
	cg.POST("/", func(c *gin.Context) {

		var requestBody struct {
			ID   primitive.ObjectID `json:"id" binding:"required"`
			Name string             `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name

		var existingRms models.RoomSpecialisation
		err := database.MongoDB.Collection("RoomSpecialisation").FindOne(c, bson.M{"Name": name}).Decode(&existingRms)

		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "RoomSpecialisation already exists"})
			fmt.Println("RoomSpecialisation already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
			return
		}

		newRoomSpecialisation := models.RoomSpecialisation{
			ID:   primitive.NewObjectID(),
			Name: name,
		}

		database.MongoDB.Collection("RoomSpecialisation").InsertOne(c, newRoomSpecialisation, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created RoomSpecialisation"})

	})

	cg.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var roomspecialisation models.RoomSpecialisation
		err = database.MongoDB.Collection("RoomSpecialisation").FindOne(c, bson.M{"_id": objectID}).Decode(&roomspecialisation)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "RoomSpecialisation not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.JSON(http.StatusOK, roomspecialisation)
	})

	cg.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var existingrms models.RoomSpecialisation
		err = database.MongoDB.Collection("RoomSpecialisation").FindOne(c, bson.M{"_id": objectID}).Decode(&existingrms)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "RoomSpecialisation not found"})
			log.Println(err)
			return
		}

		var requestBody struct {
			ID   primitive.ObjectID `json:"id"`
			Name string             `json:"name"`
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

		result, err := database.MongoDB.Collection("RoomSpecialisation").UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": "No data provided to update"})
			log.Println("Warning: No data provided to update the RoomSpecialisation")
			return
		}
		var updatedRoomSpecialisation models.RoomSpecialisation
		err = database.MongoDB.Collection("RoomSpecialisation").FindOne(c, bson.M{"_id": objectID}).Decode(&updatedRoomSpecialisation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, updatedRoomSpecialisation)
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("RoomSpecialisation").DeleteOne(c, bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "RoomSpecialisation not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "RoomSpecialisation deleted"})
	})
}
