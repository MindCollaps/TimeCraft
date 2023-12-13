package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type RoomSpecialisation struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type RoomSpecialisationStruct struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func RoomSpecialisationToStruct(c *gin.Context, roomSpecialisation RoomSpecialisation) (RoomSpecialisationStruct, error) {
	roomSpecialisationStruct := RoomSpecialisationStruct{
		ID:   roomSpecialisation.ID.Hex(),
		Name: roomSpecialisation.Name,
	}

	return roomSpecialisationStruct, nil
}

func LoadRoomSpecialisations(c *gin.Context, specialisationIDs []primitive.ObjectID) ([]RoomSpecialisationStruct, error) {
	var specialisationsStruct []RoomSpecialisationStruct
	for _, specialisationID := range specialisationIDs {
		var specialisation RoomSpecialisation
		err := database.MongoDB.Collection("RoomSpecialisations").FindOne(c, bson.M{
			"_id": specialisationID,
		}).Decode(&specialisation)

		if err != nil {
			continue
		}

		specialisationStruct, err := RoomSpecialisationToStruct(c, specialisation)
		if err != nil {
			continue
		}
		specialisationsStruct = append(specialisationsStruct, specialisationStruct)
	}
	return specialisationsStruct, nil
}
