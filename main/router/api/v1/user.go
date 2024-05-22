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
	joi "src/main/core"
	"src/main/crypt"
	"src/main/database"
	"src/main/database/models"
	"strings"
)

// /api/v1/usr/login
func userHandler(cg *gin.RouterGroup) {
	cg.POST("/login", func(c *gin.Context) {
		//check body for username and password
		var requestBody struct {
			Email    string `form:"email"  binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		email := strings.ToLower(requestBody.Email)
		password := requestBody.Password

		//check if user exists
		var user models.User
		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"email": email}).Decode(&user)
		//if user exists, check password using crypt
		if err == nil {
			if crypt.CheckPasswordHash(password, user.Password) {
				//generate jwt token
				token, err := crypt.GenerateLoginToken(user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Token generation failed"})
					log.Println(err)
					return
				}

				//set cookie with age of 2 days, setting maxAge to: 3600 * 24 * 2
				c.SetCookie("auth", token, 3600*24*2, "/", "", false, false)

				c.JSON(http.StatusOK, gin.H{"msg": "Logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "An error occurred", "error": "Invalid credentials"})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "An error occurred", "error": "Invalid credentials"})
			log.Println(err)
			return
		}
	})

	cg.POST("/register", func(c *gin.Context) {
		var requestBody struct {
			PasswordRepeat string `form:"passwordRepeat" binding:"required"`
			Password       string `form:"password" binding:"required"`
			Email          string `form:"email" binding:"required"`
		}

		if err := c.ShouldBind(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			log.Println(err)
			return
		}

		//joi validation
		if err := joi.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Password invalid. Please follow the password rules"})
			log.Println(err)
			return
		}

		if err := joi.EmailSchema.Validate(requestBody.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Email invalid. Please follow the email rules"})
			log.Println(err)
			return
		}

		if requestBody.Password != requestBody.PasswordRepeat {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Passwords do not match"})
			return
		}

		password := requestBody.Password
		email := strings.ToLower(requestBody.Email)

		// Check if the user already exists in the database by querying with the username
		var existingUser models.User

		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"email": email}).Decode(&existingUser)

		// when the email exists, there is no error --> only continue with error, because then the email doesn't exist yet
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "Username or Email already exists"})
			log.Println("Username or Email already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			// mongo.ErrNoDocuments is expected when the user does not exist
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			log.Println(err)
			return
		}

		hashedPassword, err := crypt.HashPassword(password)

		newUser := models.User{
			ID:                 primitive.NewObjectID(),
			Password:           hashedPassword,
			Email:              email,
			StaredTimeTableIds: []primitive.ObjectID{},
		}

		result, err := database.MongoDB.Collection("user").InsertOne(c, newUser, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Created user", "id": result.InsertedID})
	})

	cg.POST("/favor", func(c *gin.Context) {
		var requestBody struct {
			UserID      primitive.ObjectID `json:"userId" binding:"required"`
			TimeTableID primitive.ObjectID `json:"timeTableId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			return
		}

		if requestBody.UserID == primitive.NilObjectID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid userID"})
			return
		}

		if requestBody.TimeTableID == primitive.NilObjectID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid timeTableID"})
			return
		}

		var user bson.M
		err := database.MongoDB.Collection("user").FindOne(
			c,
			bson.M{"_id": requestBody.UserID, "staredTimeTableIds": requestBody.TimeTableID},
		).Decode(&user)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Favorite already exists"})
			return
		}

		result, err := database.MongoDB.Collection("user").UpdateOne(
			c,
			bson.M{"_id": requestBody.UserID},
			bson.M{"$addToSet": bson.M{"staredTimeTableIds": requestBody.TimeTableID}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Favorite added successfully"})
	})

	cg.DELETE("/favor", func(c *gin.Context) {
		var requestBody struct {
			UserID      primitive.ObjectID `json:"userId" binding:"required"`
			TimeTableID primitive.ObjectID `json:"timeTableId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			return
		}

		if requestBody.UserID == primitive.NilObjectID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid userID"})
			return
		}

		if requestBody.TimeTableID == primitive.NilObjectID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid timeTableID"})
			return
		}

		var user bson.M
		err := database.MongoDB.Collection("user").FindOne(
			c,
			bson.M{"_id": requestBody.UserID},
		).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "User not found"})
			return
		}

		err = database.MongoDB.Collection("user").FindOne(
			c,
			bson.M{"_id": requestBody.UserID, "staredTimeTableIds": requestBody.TimeTableID},
		).Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Favorite not found"})
			return
		}

		result, err := database.MongoDB.Collection("user").UpdateOne(
			c,
			bson.M{"_id": requestBody.UserID},
			bson.M{"$pull": bson.M{"staredTimeTableIds": requestBody.TimeTableID}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database error"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Favorite removed successfully"})
	})

	cg.GET("/favor/:userId", func(c *gin.Context) {
		userID := c.Param("userId")

		userIDObj, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid userID"})
			return
		}

		var user struct {
			StaredTimeTableIds []primitive.ObjectID `bson:"staredTimeTableIds"`
		}

		err = database.MongoDB.Collection("user").FindOne(
			c,
			bson.M{"_id": userIDObj},
		).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"staredTimeTableIds": user.StaredTimeTableIds})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("user").DeleteOne(c, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			log.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "An error occurred", "error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Deleted user"})
	})
}
