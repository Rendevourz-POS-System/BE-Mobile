package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

type Shelter struct {
	ID                   primitive.ObjectID   `json:"Id" bson:"_id,omitempty"`
	UserId               primitive.ObjectID   `json:"UserId" bson:"user_id"`
	ShelterLocation      primitive.ObjectID   `json:"ShelterLocation" bson:"shelter_location" validate:"required"`
	ShelterName          string               `json:"ShelterName" bson:"shelter_name" validate:"required"`
	ShelterAddress       string               `json:"ShelterAddress" bson:"shelter_address" validate:"required"`
	ShelterCapacity      int                  `json:"ShelterCapacity" bson:"shelter_capacity" validate:"required,number"`
	ShelterContactNumber string               `json:"ShelterContactNumber" bson:"shelter_contact_number" validate:"required,min=10"`
	ShelterDescription   string               `json:"ShelterDescription,omitempty" bson:"shelter_description" default:""`
	TotalPet             int                  `json:"TotalPet" bson:"total_pet" default:"0"`
	BankAccountNumber    string               `json:"BankAccountNumber" bson:"bank_account_number" validate:"omitempty,required,min=10"`
	PetTypeAccepted      []primitive.ObjectID `json:"PetTypeAccepted" bson:"pet_type_accepted" validate:"required,pet-accepted-min"`
	Image                []string             `json:"ImagePath" bson:"image" validate:"omitempty"`
	Pin                  string               `json:"Pin" bson:"pin" validate:"omitempty,required,min=6,max=8"`
	ShelterVerified      bool                 `json:"ShelterVerified" bson:"shelter_verified" default:"false"`
	OldImage             []string             `json:"OldImagePath,omitempty"`
	CreatedAt            *time.Time           `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	UpdatedAt            *time.Time           `json:"UpdatedAt" bson:"UpdatedAt,omitempty"`
	DeletedAt            *time.Time           `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	// ShelterSearch
	ShelterSearch struct {
		Search              string             `json:"Search"`
		Page                int                `json:"Page"`
		PageSize            int                `json:"PageSize"`
		OrderBy             string             `json:"OrderBy"`
		Sort                string             `json:"Sort"`
		ShelterLocationName string             `json:"ShelterLocationName"`
		PetType             string             `json:"PetType"`
		ShelterId           primitive.ObjectID `json:"ShelterId"`
		SearchByUserId      primitive.ObjectID `json:"SearchByUserId"`
		UserId              primitive.ObjectID `json:"UserId"`
	}

	// ShelterCreate
	ShelterCreate struct {
		Files   *multipart.FileHeader `form:"Files" bson:"-" validate:"omitempty"`
		Shelter Shelter               `form:"Shelter" bson:"Shelter" validate:"required"`
	}

	// SehlterUpdate
	ShelterUpdate struct {
		Shelter *Shelter `form:"Shelter" bson:"Shelter" validate:"required"`
	}
	// ShelterResponsePayload
	ShelterResponsePayload struct {
		ID                   primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
		UserId               primitive.ObjectID `json:"UserId" bson:"user_id"`
		ShelterLocation      primitive.ObjectID `json:"ShelterLocation" bson:"shelter_location"`
		ShelterLocationName  string             `json:"ShelterLocationName" bson:"shelter_location_name"`
		ShelterName          string             `json:"ShelterName" bson:"shelter_name"`
		ShelterAddress       string             `json:"ShelterAddress" bson:"shelter_address"`
		ShelterCapacity      int                `json:"ShelterCapacity" bson:"shelter_capacity"`
		ShelterContactNumber string             `json:"ShelterContactNumber" bson:"shelter_contact_number"`
		ShelterDescription   string             `json:"ShelterDescription,omitempty" bson:"shelter_description"`
		TotalPet             int                `json:"TotalPet" bson:"total_pet" default:"0"`
		BankAccountNumber    string             `json:"BankAccountNumber" bson:"bank_account_number"`
		PetTypeAccepted      []string           `json:"PetTypeAccepted" bson:"pet_type_accepted"`
		PetTypeAcceptedName  []string           `json:"PetTypeAcceptedName" bson:"pet_type_accepted_name"`
		Image                []string           `json:"ImagePath" bson:"image"`
		ImageBase64          []string           `json:"ImageBase64"`
		Pin                  string             `json:"Pin" bson:"pin"`
		ShelterVerified      bool               `json:"ShelterVerified" bson:"shelter_verified"`
		CreatedAt            *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
		DeletedAt            *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
	}
)
