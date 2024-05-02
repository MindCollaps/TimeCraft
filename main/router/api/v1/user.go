package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	joi "src/main/core"
	"src/main/crypt"
	"src/main/database"
	"src/main/database/models"
	"strings"
)

func userHandler(cg *gin.RouterGroup) {
	//    /api/v1/usr/login
	cg.POST("/login", func(c *gin.Context) {
		//check body for username and password
		var requestBody struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password

		//check if user exists
		var user models.User
		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"username": username}).Decode(&user)
		//if user exists, check password using crypt
		if err == nil {
			if crypt.CheckPasswordHash(password, user.Password) {
				//generate jwt token
				token, err := crypt.GenerateLoginToken(user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
					return
				}

				//set cookie with age of 2 days, setting maxAge to: 3600 * 24 * 2
				c.SetCookie("auth", token, 3600*24*2, "/", "", false, false)

				c.JSON(http.StatusOK, gin.H{"msg": "Logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "Not Authorized"})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Not Authorized"})
		}
	})

	cg.POST("/register", func(c *gin.Context) {
		var requestBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			fmt.Println(err)
			return
		}

		//joi validation
		if err := joi.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Username invalid", "field": "username"})
			fmt.Println(err)
			return
		}

		if err := joi.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Password invalid", "field": "password"})
			fmt.Println(err)
			return
		}

		if err := joi.EmailSchema.Validate(requestBody.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Email invalid", "field": "email"})
			fmt.Println(err)
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password
		email := strings.ToLower(requestBody.Email)

		// Check if the user already exists in the database by querying with the username
		var existingUser models.User
		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"username": username}).Decode(&existingUser)

		if err == nil {
			// User with the same username already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "Username already exists"})
			fmt.Println("Username already exists")
			return
		}

		err = database.MongoDB.Collection("user").FindOne(c, bson.M{"email": email}).Decode(&existingUser)

		if err == nil {
			// User with the same email already exists
			c.JSON(http.StatusConflict, gin.H{"msg": "Email already exists"})
			fmt.Println("Email already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			// Handle other database query errors
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
			return
		}

		hashedPassword, err := crypt.HashPassword(password)

		newUser := models.User{
			ID:       primitive.NewObjectID(),
			Username: username,
			Password: hashedPassword,
			Email:    email,
		}

		_, err = database.MongoDB.Collection("user").InsertOne(c, newUser, options.InsertOne())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Created user"})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID"})
			return
		}
		result, err := database.MongoDB.Collection("user").DeleteOne(c, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			fmt.Println(err)
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Deleted user"})
	})
}
