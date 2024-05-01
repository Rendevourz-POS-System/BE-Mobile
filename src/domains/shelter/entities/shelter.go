package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Shelter struct {
	ID                   string             `json:"Id" bson:"_id,omitempty"`
	UserId               primitive.ObjectID `json:"UserId" bson:"user_id"`
	ShelterName          string             `json:"ShelterName" bson:"shelter_name" validate:"required"`
	ShelterLocation      string             `json:"ShelterLocation" bson:"shelter_location" validate:"required"`
	ShelterAddress       string             `json:"ShelterAddress" bson:"shelter_address" validate:"required"`
	ShelterCapacity      int                `json:"ShelterCapacity" bson:"shelter_capacity" validate:"required,number"`
	ShelterContactNumber string             `json:"ShelterContactNumber" bson:"shelter_contact_number" validate:"required,min=10"`
	ShelterDescription   string             `json:"ShelterDescription,omitempty" bson:"shelter_description" default:""`
	TotalPet             int                `json:"TotalPet" bson:"total_pet" default:"0"`
	BankAccountNumber    string             `json:"BankAccountNumber" bson:"bank_account_number" validate:"omitempty,required,min=10"`
	PetTypeAccepted      []string           `json:"PetTypeAccepted" bson:"pet_type_accepted" validate:"required,pet-accepted-min"`
	ImagePath            []string           `json:"ImagePath" bson:"image" validate:"omitempty"`
	Pin                  string             `json:"Pin" bson:"pin" validate:"omitempty,required,min=6,max=8"`
	ShelterVerified      bool               `json:"ShelterVerified" bson:"shelter_verified" default:"false"`
	CreatedAt            *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	DeletedAt            *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	ShelterSearch struct {
		Search    string             `json:"Search"`
		Page      int                `json:"Page"`
		PageSize  int                `json:"PageSize"`
		OrderBy   string             `json:"OrderBy"`
		Sort      string             `json:"Sort"`
		ShelterId primitive.ObjectID `json:"ShelterId"`
		UserId    primitive.ObjectID `json:"UserId"`
	}
)
