package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PetType struct {
	ID   primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	Type string             `json:"Type" bson:"type" validate:"required"`
}
