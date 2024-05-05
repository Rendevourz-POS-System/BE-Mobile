package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShelterLocation struct {
	ID           primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	LocationName string             `json:"LocationName" bson:"location_name" binding:"required"`
}
