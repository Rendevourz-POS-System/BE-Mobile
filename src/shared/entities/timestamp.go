package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Timestamp struct {
	CreatedAt primitive.Timestamp `json:"created_at"`
	UpdateAt  primitive.Timestamp `json:"updated_at"`
	DeletedAt primitive.Timestamp `json:"deleted_at"`
}
