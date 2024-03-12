package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	Email    string             `json:"Email" bson:"email" validate:"required,email"`
	Username string             `json:"Username" bson:"username" validate:"required,min=4"`
	Password string             `json:"Password" bson:"password" validate:"required,min=8,alphanum_symbol"`
	Role     string             `json:"Role" bson:"role"`
	IsActive bool               `json:"IsActive" bson:"is_active"`
}
