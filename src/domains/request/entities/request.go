package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Request struct {
	Id          primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"UserId" bson:"user_id" validate:"required"`
	ShelterId   primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
	Type        string             `json:"TypeId" bson:"type_id" validate:"required"`
	Status      string             `json:"StatusId" bson:"status_id"`
	Reason      *string            `json:"Reason,omitempty" bson:"reason"`
	RequestedAt *time.Time         `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
	CompletedAt *time.Time         `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
}
