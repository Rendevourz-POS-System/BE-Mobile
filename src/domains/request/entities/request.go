package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/domains/request/presistence"
	User "main.go/domains/user/entities"
	"time"
)

type Request struct {
	Id        primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"UserId" bson:"user_id"`
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
		UserId    primitive.ObjectID `json:"UserId" bson:"user_id"`
		ShelterId primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
		RequestId primitive.ObjectID `json:"RequestId,omitempty" bson:"request_id"`
		Type      string             `json:"Type" bson:"type" validate:"required,donations"`
		Status    presistence.Status `json:"Status" bson:"status_id" validate:"omitempty" default:"New"`
		//Job         string             `json:"Job" bson:"job"`
		Reason      *string    `json:"Reason,omitempty" bson:"reason"`
		Amount      int64      `json:"Amount" validate:"omitempty"`
		PaymentType string     `json:"PaymentType" bson:"payment_type" validate:"required,payment_type"`
		BankType    *string    `json:"BankType,omitempty" bson:"bank_type" validate:"omitempty,bank_type"`
		EWalletType *string    `json:"EWalletType,omitempty" bson:"e_wallet_type" validate:"omitempty,e_wallet"`
		RequestedAt *time.Time `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
		CompletedAt *time.Time `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
	}
	RequestResponse struct {
		Request         *Request         `json:"Request,omitempty"`
		Donation        *DonationShelter `json:"Donation,omitempty"`
		Adoption        *AdoptionShelter `json:"Adoption,omitempty"`
		DonationPayload *DonationPayload `json:"DonationPayload,omitempty"`
		User            *User.User       `json:"User,omitempty"`
		UserTarget      *User.User       `json:"UserTarget,omitempty"`
	}
)
