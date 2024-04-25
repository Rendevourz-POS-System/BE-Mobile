package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShelterFavorite struct {
	ShelterId primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
	UserId    primitive.ObjectID `json:"UserId" bson:"user_id" validate:"required"`
}

type (
	// ShelterFavoriteCreate Payload
	ShelterFavoriteCreate struct {
		UserId    primitive.ObjectID `json:"UserId"`
		ShelterId primitive.ObjectID `json:"ShelterId" validate:"required"`
	}
)
