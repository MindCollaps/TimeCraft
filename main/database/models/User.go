package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type User struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Email              string               `json:"email" bson:"email"`
	Password           string               `json:"password" bson:"password"`
	IsVerifiedEmail    bool                 `json:"isVerifiedEmail" bson:"isVerifiedEmail"`
	Username           string               `json:"username" bson:"username"`
	Permissions        []primitive.ObjectID `json:"permissions" bson:"permissions"`
	SSOToken           string               `json:"SSOToken" bson:"SSOToken"`
	StaredTimeTableIds []primitive.ObjectID `json:"staredTimeTableIds" bson:"staredTimeTableIds"`
}

type UserStruct struct {
	ID                 string             `json:"id" bson:"_id"`
	Email              string             `json:"email" bson:"email"`
	Password           string             `json:"password" bson:"password"`
	IsVerifiedEmail    bool               `json:"isVerifiedEmail" bson:"isVerifiedEmail"`
	Username           string             `json:"username" bson:"username"`
	Permissions        []PermissionStruct `json:"permissions" bson:"permissions"`
	SSOToken           string             `json:"SSOToken" bson:"SSOToken"`
	StaredTimeTableIds []TimeTableStruct  `json:"staredTimeTableIds" bson:"staredTimeTableIds"`
}

func UserToStruct(c *gin.Context, user User) (UserStruct, error) {
	userStruct := UserStruct{
		ID:              user.ID.Hex(),
		Email:           user.Email,
		Password:        user.Password,
		IsVerifiedEmail: user.IsVerifiedEmail,
		Username:        user.Username,
		SSOToken:        user.SSOToken,
	}

	// Load Permissions
	permissions, err := LoadPermissions(c, user.Permissions)
	if err != nil {
		return UserStruct{}, err
	}
	userStruct.Permissions = permissions

	// Load StaredTimeTableIds
	staredTimeTables, err := LoadTimeTables(c, user.StaredTimeTableIds)
	if err != nil {
		return UserStruct{}, err
	}
	userStruct.StaredTimeTableIds = staredTimeTables

	return userStruct, nil
}

func LoadUsers(c *gin.Context, userIDs []primitive.ObjectID) ([]UserStruct, error) {
	var usersStruct []UserStruct
	for _, userID := range userIDs {
		var user User
		err := database.MongoDB.Collection("Users").FindOne(c, bson.M{
			"_id": userID,
		}).Decode(&user)

		if err != nil {
			continue
		}

		userStruct, err := UserToStruct(c, user)
		if err != nil {
			continue
		}
		usersStruct = append(usersStruct, userStruct)
	}
	return usersStruct, nil
}
