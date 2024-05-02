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

func tblHandler(cg *gin.RouterGroup) {
	//    /api/v1/tbl/...
	cg.POST("/create", func(c *gin.Context) {
		//check body for username and password
		var requestBody struct {
			Id   primitive.ObjectID   `json:"id" binding:"required"`
			Name string               `json:"name" binding:"required"`
			Days []primitive.ObjectID `json:"days" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		name := requestBody.Name
		days := requestBody.Days

		var existingTbl models.TimeTable
		err := database.MongoDB.Collection("TimeTable").FindOne(c, bson.M{"Name": name}).Decode(&existingTbl)

		if err == nil {
			// Table with the same name already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "TimeTable already exists"})
			log.Println("TimeTable already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			log.Println(err)
			return
		}

		newTable := models.TimeTable{
			ID:   primitive.NewObjectID(),
			Name: name,
			Days: days,
		}

		database.MongoDB.Collection("TimeTable").InsertOne(c, newTable, options.InsertOne())
		c.JSON(http.StatusOK, gin.H{"msg": "Created Timetable"})

	})

	cg.GET("/{id}", func(c *gin.Context) {

		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Please give correct ID"})
			return
		} else {
			//Abfrage aus der DB
		}
	})

	cg.PATCH("/{id}", func(c *gin.Context) {

		//Updatefunktion
	})

	cg.DELETE("/{id}", func(c *gin.Context) {

		//Deletefunktion
	})
}
