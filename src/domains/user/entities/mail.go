package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type VerifiedEmail struct {
	Id         primitive.ObjectID `json:"Id" bson:"_id"`
	UserId     primitive.ObjectID `json:"UserId" bson:"UserId"`
	SecretCode string             `json:"SecretCode" bson:"SecretCode"`
	IsUsed     bool               `json:"IsUsed" bson:"IsUsed"`
	CreatedAt  time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	ExpiredAt  time.Time          `json:"ExpiredAt" bson:"ExpiredAt"`
}

type (
	// MailSend Payload for sending mail (Send To Email)
	MailSend struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Content string `json:"content"`
		Cc      string `json:"cc"`
		Bcc     string `json:"bcc"`
		Attach  string `json:"attach"`
	}
	// MailVerifyResponse Response for mail verification
	MailVerifyResponse struct {
		Email   string `json:"email"`
		Message string `json:"message"`
	}
)
