package entities

import "time"

type VerifiedEmail struct {
	Id         uint      `json:"Id" bson:"_id"`
	UserId     uint      `json:"UserId" bson:"UserId"`
	SecretCode string    `json:"SecretCode" bson:"SecretCode"`
	IsUsed     bool      `json:"IsUsed" bson:"IsUsed"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"CreatedAt" bson:"CreatedAt"`
	ExpiredAt  time.Time `gorm:"autoCreateTime + interval '15 minutes'" json:"ExpiredAt" bson:"ExpiredAt"`
}

type GmailSender struct {
	Name              string
	FromEmailAddress  string
	FromEmailPassword string
}
