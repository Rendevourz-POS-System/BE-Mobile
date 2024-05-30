package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DonationShelter struct {
	Id                primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	RequestId         primitive.ObjectID `json:"RequestId" bson:"request_id"`
	Amount            int64              `json:"Amount" bson:"amount"`
	TransactionDate   *time.Time         `json:"TransactionDate" bson:"transactionDate" validate:"required"`
	CreatedAt         *time.Time         `json:"CreatedAt" bson:"CreatedAt" validate:"required"`
	StatusTransaction string             `json:"StatusTransaction" bson:"status_transactionDate" default:"Ongoing"`
	PaymentType       string             `json:"PaymentType" bson:"payment_type" validate:"required"`
	Channel           string             `json:"Channel" bson:"channel" `
}
