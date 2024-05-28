package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DonationShelter struct {
	Id                primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	RequestId         primitive.ObjectID `json:"RequestId" bson:"request_id"`
	Amount            float64            `json:"Amount" bson:"amount"`
	TransactionDate   time.Time          `json:"TransactionDate" bson:"transactionDate" validate:"required"`
	StatusTransaction string             `json:"StatusTransaction" bson:"status_transactionDate" default:"Ongoing"`
	PaymentType       string             `json:"PaymentType" bson:"payment_type" validate:"required"`
}
