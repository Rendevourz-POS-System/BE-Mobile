package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type PetFavorite struct {
	PetId  primitive.ObjectID `json:"PetId" bson:"pet_id" validate:"required"`
	UserId primitive.ObjectID `json:"UserId" bson:"user_id" validate:"required"`
}

type (
	// PetFavoriteCreate Payload
	PetFavoriteCreate struct {
		UserId primitive.ObjectID `json:"UserId"`
		PetId  primitive.ObjectID `json:"PetId" validate:"required"`
	}
)
