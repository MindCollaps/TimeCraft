package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID                 string       `json:"id" bson:"_id"`
	Email              string       `json:"email" bson:"email"`
	Password           string       `json:"password" bson:"password"`
	IsVerifiedEmail    bool         `json:"isVerifiedEmail" bson:"isVerifiedEmail"`
	Username           string       `json:"username" bson:"username"`
	Permissions        []Permission `json:"permissions" bson:"permissions"`
	SSOToken           string       `json:"SSOToken" bson:"SSOToken"`
	StaredTimeTableIds []TimeTable  `json:"staredTimeTableIds" bson:"staredTimeTableIds"`
}
