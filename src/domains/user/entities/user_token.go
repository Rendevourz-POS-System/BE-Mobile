package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserToken struct {
	Id        primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"UserId" bson:"UserId"`
	Token     string             `json:"Token" bson:"Token"`
	IsUsed    bool               `json:"IsUsed" bson:"IsUsed"`
	CreatedAt *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	ExpiredAt *time.Time         `json:"ExpiredAt" bson:"ExpiredAt,omitempty"`
	DeletedAt *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}
