package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

type Pet struct {
	ID             primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	ShelterId      primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
	PetName        string             `json:"PetName" bson:"pet_name" validate:"required"`
	PetType        string             `json:"PetType" bson:"pet_type" validate:"required"`
	PetAge         int                `json:"PetAge" bson:"pet_age" validate:"required,number,pet-age"`
	PetGender      string             `json:"PetGender" bson:"pet_gender" validate:"omitempty,required,pet-gender"`
	PetStatus      bool               `json:"PetStatus" bson:"pet_status" validate:"omitempty" default:"false"`
	PetDescription string             `json:"PetDescription" bson:"pet_description" validate:"omitempty,required,min=10"`
	IsVaccinated   bool               `json:"IsVaccinated" bson:"is_vaccinated" validate:"omitempty,required"`
	ImagePath      []string           `json:"Image" bson:"image" validate:"omitempty"`
	ImageBase64    []string           `json:"ImageBase64" validate:"omitempty"`
	PetDob         *time.Time         `json:"PetDob" bson:"pet_dob" validate:"omitempty"`
	CreatedAt      *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	DeletedAt      *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	// Pet Response Payload
	PetResponsePayload struct {
		ID              primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
		ShelterId       primitive.ObjectID `json:"ShelterId" bson:"shelter_id"`
		ShelterName     string             `json:"ShelterName" bson:"shelter_name"`
		ShelterLocation string             `json:"ShelterLocation" bson:"shelter_location"`
		Location        string             `json:"Location" bson:"shelter_location_name"`
		PetName         string             `json:"PetName" bson:"pet_name"`
		PetType         string             `json:"PetType" bson:"pet_type"`
		PetGender       string             `json:"PetGender" bson:"pet_gender"`
		PetStatus       bool               `json:"PetStatus" bson:"pet_status"`
		PetDescription  string             `json:"PetDescription" bson:"pet_description"`
		IsVaccinated    bool               `json:"IsVaccinated" bson:"is_vaccinated"`
		ImagePath       []string           `json:"Image" bson:"image"`
		ImageBase64     []string           `json:"ImageBase64"`
		PetAge          int                `json:"PetAge" bson:"pet_age"`
		CreatedAt       *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
		DeletedAt       *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
	}
	// PetCreate Payload
	PetCreate struct {
		Files *multipart.FileHeader `form:"Files" bson:"-" validate:"omitempty"`
		Pet   Pet                   `form:"Pet" bson:"Pet" validate:"required"`
	}
	// PetSearch struct
	PetSearch struct {
		Search      string             `json:"Search"`
		Page        int                `json:"Page"`
		PageSize    int                `json:"PageSize"`
		OrderBy     string             `json:"OrderBy"`
		Sort        string             `json:"Sort"`
		ShelterId   primitive.ObjectID `json:"ShelterId"`
		ShelterName string             `json:"ShelterName"`
		Location    string             `json:"Location"`
		Gender      string             `json:"Gender"`
		AgeStart    int                `json:"AgeStart"`
		AgeEnd      int                `json:"AgeEnd"`
		Type        string             `json:"Type"`
	}
)
