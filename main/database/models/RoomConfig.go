package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type RoomConfig struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id"`
	Name           string               `json:"name" bson:"name"`
	RoomNr         string               `json:"roomNr" bson:"roomNr"` // All room number have 2 digits --> 1.0, 3.1, ....
	Capacity       int                  `json:"capacity" bson:"capacity"`
	ExamCapacity   int                  `json:"examCapacity" bson:"examCapacity"`
	Blocks         []primitive.ObjectID `json:"blocks" bson:"blocks"`
	Specialisation []primitive.ObjectID `json:"specialisation" bson:"specialisation"`
}

type RoomConfigStruct struct {
	ID             string                     `json:"id" bson:"_id"`
	Name           string                     `json:"name" bson:"name"`
	RoomNr         string                     `json:"roomNr" bson:"roomNr"` // All room number have 2 digits --> 1.0, 3.1, ....
	Capacity       int                        `json:"capacity" bson:"capacity"`
	ExamCapacity   int                        `json:"examCapacity" bson:"examCapacity"`
	Blocks         []RoomConfigStruct         `json:"blocks" bson:"blocks"`
	Specialisation []RoomSpecialisationStruct `json:"specialisation" bson:"specialisation"`
}

func RoomConfigToStruct(c *gin.Context, roomConfig RoomConfig) (RoomConfigStruct, error) {
	roomConfigStruct := RoomConfigStruct{
		ID:           roomConfig.ID.Hex(),
		Name:         roomConfig.Name,
		RoomNr:       roomConfig.RoomNr,
		Capacity:     roomConfig.Capacity,
		ExamCapacity: roomConfig.ExamCapacity,
	}

	// Convert Blocks to []RoomConfigStruct
	blocks, err := LoadRoomConfigs(c, roomConfig.Blocks)
	if err != nil {
		return RoomConfigStruct{}, err
	}

	roomConfigStruct.Blocks = blocks

	// Convert Specialisation to []RoomSpecialisationStruct
	specialisation, err := LoadRoomSpecialisations(c, roomConfig.Specialisation)
	if err != nil {
		return RoomConfigStruct{}, err
	}

	roomConfigStruct.Specialisation = specialisation

	return roomConfigStruct, nil
}

func LoadRoomConfigs(c *gin.Context, roomConfigIDs []primitive.ObjectID) ([]RoomConfigStruct, error) {
	var roomConfigsStruct []RoomConfigStruct
	for _, roomConfigID := range roomConfigIDs {
		var roomConfig RoomConfig
		err := database.MongoDB.Collection("RoomConfigs").FindOne(c, bson.M{
			"_id": roomConfigID,
		}).Decode(&roomConfig)

		if err != nil {
			continue
		}

		roomConfigStruct, err := RoomConfigToStruct(c, roomConfig)
		if err != nil {
			continue
		}
		roomConfigsStruct = append(roomConfigsStruct, roomConfigStruct)
	}
	return roomConfigsStruct, nil
}
