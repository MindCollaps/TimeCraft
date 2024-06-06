package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ContainsNilObjectID(array []primitive.ObjectID) bool {
	for _, id := range array {
		if id.IsZero() || id == primitive.NilObjectID {
			return true
		}
	}
	return false
}
