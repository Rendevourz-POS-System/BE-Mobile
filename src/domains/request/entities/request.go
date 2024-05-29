package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/domains/request/presistence"
	"time"
)

type Request struct {
	Id        primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"UserId" bson:"user_id" validate:"required"`
	ShelterId primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
	Type      string             `json:"Type" bson:"type" validate:"required,request-type"`
	Status    presistence.Status `json:"Status" bson:"status_id" validate:"omitempty" default:"New"`
	//Job         string             `json:"Job" bson:"job"`
	Reason      *string    `json:"Reason,omitempty" bson:"reason"`
	RequestedAt *time.Time `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
	CompletedAt *time.Time `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
}

type (
	RescuePayload struct {
		Request `json:"Request" validate:"required"`
	}
	AdoptionPayload struct {
		Request `json:"Request" validate:"required"`
	}
	DonationPayload struct {
		Id        primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
		UserId    primitive.ObjectID `json:"UserId" bson:"user_id" validate:"required"`
		ShelterId primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
		Type      string             `json:"Type" bson:"type" validate:"required,request-type"`
		Status    presistence.Status `json:"Status" bson:"status_id" validate:"omitempty" default:"New"`
		//Job         string             `json:"Job" bson:"job"`
		Reason      *string    `json:"Reason,omitempty" bson:"reason"`
		Amount      float64    `json:"Amount" validate:"omitempty"`
		PaymentType string     `json:"PaymentType" bson:"payment_type" validate:"required,payment_type"`
		RequestedAt *time.Time `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
		CompletedAt *time.Time `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
	}
)
