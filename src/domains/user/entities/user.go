package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"Username" bson:"Username"`
	Password string             `json:"Password" bson:"Password"`
	Email    string             `json:"Email" bson:"Email"`
	Role     string             `json:"Role" bson:"Role"`
}
