package v1

import (
	"errors"
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
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			fmt.Println(err)
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
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Token generation failed"})
					fmt.Println(err)
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
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "An error occurred", "error": "Database error"})
			fmt.Println(err)
			return
		}
	})

	cg.POST("/register", func(c *gin.Context) {
		var requestBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid request body"})
			fmt.Println(err)
			return
		}

		//joi validation
		if err := joi.UsernameSchema.Validate(requestBody.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Username invalid. Please follow the username rules"})
			fmt.Println(err)
			return
		}

		if err := joi.PasswordSchema.Validate(requestBody.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Password invalid. Please follow the password rules"})
			fmt.Println(err)
			return
		}

		if err := joi.EmailSchema.Validate(requestBody.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Email invalid. Please follow the email rules"})
			fmt.Println(err)
			return
		}

		username := strings.ToLower(requestBody.Username)
		password := requestBody.Password
		email := strings.ToLower(requestBody.Email)

		// Check if the user already exists in the database by querying with the username
		var existingUser models.User
		err := database.MongoDB.Collection("user").FindOne(c, bson.M{"username": username}).Decode(&existingUser)

		// when the user exists, there is no error --> only continue with error, because then the user doesn't exist yet
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "Username or Email already exists"})
			fmt.Println("Username or Email already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			// mongo.ErrNoDocuments is expected when the user does not exist
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			fmt.Println(err)
			return
		}

		err = database.MongoDB.Collection("user").FindOne(c, bson.M{"email": email}).Decode(&existingUser)

		// when the email exists, there is no error --> only continue with error, because then the email doesn't exist yet
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"msg": "An error occurred", "error": "Username or Email already exists"})
			fmt.Println("Username or Email already exists")
			return
		} else if !errors.Is(err, mongo.ErrNoDocuments) {
			// Handle other database query errors
			// mongo.ErrNoDocuments is expected when the user does not exist
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Created user"})
	})

	cg.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "An error occurred", "error": "Invalid ID"})
			return
		}
		_, err = database.MongoDB.Collection("user").DeleteOne(c, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "An error occurred", "error": "Database	error"})
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Deleted user"})
	})
}
