package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pet struct {
	ID             primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	ShelterId      primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
	PetName        string             `json:"PetName" bson:"pet_name" validate:"required"`
	PetType        string             `json:"PetType" bson:"pet_type" validate:"required"`
	PetAge         int                `json:"PetAge" bson:"pet_age" validate:"required,number"`
	PetGender      string             `json:"PetGender" bson:"pet_gender" validate:"required,pet-gender"`
	PetStatus      string             `json:"PetStatus" bson:"pet_status" validate:"required"`
	PetDescription string             `json:"PetDescription" bson:"pet_description" validate:"omitempty,required,min=10"`
	Image          string             `json:"Image" bson:"image" validate:"omitempty"`
	CreatedAt      *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	DeletedAt      *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	PetSearch struct {
		Search    string             `json:"Search"`
		Page      int                `json:"Page"`
		PageSize  int                `json:"PageSize"`
		OrderBy   string             `json:"OrderBy"`
		Sort      string             `json:"Sort"`
		ShelterId primitive.ObjectID `json:"ShelterId"`
		Location  string             `json:"Location"`
		Gender    string             `json:"Gender"`
		AgeStart  int                `json:"AgeStart"`
		AgeEnd    int                `json:"AgeEnd"`
		Type      string             `json:"Type"`
	}
)
