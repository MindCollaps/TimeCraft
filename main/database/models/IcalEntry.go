package models

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/main/database"
)

type IcalEntry struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Text        string             `json:"text" bson:"text"`
	LastUpdated primitive.DateTime `json:"lastUpdated" bson:"lastUpdated"`
}

type IcalEntryStruct struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Text        string `json:"text" bson:"text"`
	LastUpdated string `json:"lastUpdated" bson:"lastUpdated"`
}

// Function to convert IcalEntry to IcalEntryStruct
func IcalEntryToStruct(icalEntry IcalEntry) IcalEntryStruct {
	icalEntryStruct := IcalEntryStruct{
		ID:          icalEntry.ID.Hex(),
		Name:        icalEntry.Name,
		Text:        icalEntry.Text,
		LastUpdated: icalEntry.LastUpdated.Time().String(),
	}
	return icalEntryStruct
}

// Function to load multiple icalEntries and convert them to IcalEntryStruct
func LoadIcalEntries(c *gin.Context, icalEntryIDs []primitive.ObjectID) ([]IcalEntryStruct, error) {
	var icalEntriesStruct []IcalEntryStruct
	for _, icalEntryID := range icalEntryIDs {
		var icalEntry IcalEntry
		err := database.MongoDB.Collection("IcalEntry").FindOne(c, bson.M{"_id": icalEntryID}).Decode(&icalEntry)
		if err != nil {
			continue
		}

		icalEntryStruct := IcalEntryToStruct(icalEntry)
		icalEntriesStruct = append(icalEntriesStruct, icalEntryStruct)
	}
	return icalEntriesStruct, nil
}
