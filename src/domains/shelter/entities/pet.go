package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

type Pet struct {
	ID             primitive.ObjectID  `json:"Id" bson:"_id,omitempty"`
	ShelterId      *primitive.ObjectID `json:"ShelterId,omitempty" bson:"shelter_id" validate:"omitempty"`
	PetName        string              `json:"PetName" bson:"pet_name" validate:"required"`
	PetType        string              `json:"PetType" bson:"pet_type" validate:"required"`
	PetAge         int                 `json:"PetAge" bson:"pet_age" validate:"required,number,pet-age"`
	PetGender      string              `json:"PetGender" bson:"pet_gender" validate:"omitempty,required,pet-gender"`
	IsAdopted      *bool               `json:"IsAdopted" bson:"is_adopted" validate:"omitempty" default:"false"`
	ReadyToAdopt   *bool               `json:"ReadyToAdopt" bson:"ready_to_adopt" validate:"omitempty" default:"false"`
	PetDescription string              `json:"PetDescription" bson:"pet_description" validate:"omitempty,required,min=10"`
	IsVaccinated   bool                `json:"IsVaccinated" bson:"is_vaccinated" validate:"omitempty,required"`
	OldImage       []string            `json:"OldImage,omitempty"`
	Image          []string            `json:"Image" bson:"image" validate:"omitempty"`
	ImageBase64    []string            `json:"ImageBase64" validate:"omitempty"`
	PetDob         *time.Time          `json:"PetDob" bson:"pet_dob" validate:"omitempty"`
	CreatedAt      *time.Time          `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	UpdatedAt      *time.Time          `json:"UpdatedAt,omitempty" bson:"UpdatedAt,omitempty"`
	DeletedAt      *time.Time          `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	// PetSearch struct
	PetSearch struct {
		Search             string             `json:"Search"`
		Page               int                `json:"Page"`
		PageSize           int                `json:"PageSize"`
		OrderBy            string             `json:"OrderBy"`
		Sort               string             `json:"Sort"`
		ShelterId          string             `json:"ShelterId,omitempty"`
		ShelterName        string             `json:"ShelterName"`
		Location           string             `json:"Location"`
		Gender             string             `json:"Gender"`
		AgeStart           int                `json:"AgeStart"`
		AgeEnd             int                `json:"AgeEnd"`
		Type               []string           `json:"Type"`
		UserId             primitive.ObjectID `json:"UserId,omitempty"`
		ReadyForAdoption   *bool              `json:"ReadyForAdoption,omitempty"`
		IsAdopted          *bool              `json:"IsAdopted,omitempty"`
		ShowPetWithShelter *bool              `json:"ShowPetWithShelter,omitempty"`
	}

	// PetUpdate Payload
	PetUpdatePayload struct {
		Files *multipart.FileHeader `form:"Files" bson:"-" validate:"omitempty"`
		Pet   Pet                   `form:"Pet" bson:"Pet" validate:"required"`
	}
	// PetCreate Payload
	PetCreate struct {
		Files *multipart.FileHeader `form:"Files" bson:"-" validate:"omitempty"`
		Pet   Pet                   `form:"Pet" bson:"Pet" validate:"required"`
	}

	// Pet Response Payload
	PetResponsePayload struct {
		ID              primitive.ObjectID  `json:"Id" bson:"_id,omitempty"`
		ShelterId       *primitive.ObjectID `json:"ShelterId" bson:"shelter_id"`
		ShelterName     string              `json:"ShelterName" bson:"shelter_name"`
		ShelterLocation string              `json:"ShelterLocation" bson:"shelter_location"`
		Location        string              `json:"Location" bson:"shelter_location_name"`
		PetName         string              `json:"PetName" bson:"pet_name"`
		PetType         string              `json:"PetType" bson:"pet_type"`
		PetGender       string              `json:"PetGender" bson:"pet_gender"`
		PetStatus       bool                `json:"PetStatus" bson:"pet_status"`
		PetDescription  string              `json:"PetDescription" bson:"pet_description"`
		IsVaccinated    bool                `json:"IsVaccinated" bson:"is_vaccinated"`
		IsAdopted       *bool               `json:"IsAdopted" bson:"is_adopted"`
		ReadyToAdopt    *bool               `json:"ReadyToAdopt" bson:"ready_to_adopt"`
		Image           []string            `json:"Image" bson:"image"`
		ImageBase64     []string            `json:"ImageBase64"`
		PetAge          int                 `json:"PetAge" bson:"pet_age"`
		CreatedAt       *time.Time          `json:"CreatedAt" bson:"CreatedAt,omitempty"`
		DeletedAt       *time.Time          `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
	}

	// Pet Delete Payload
	PetDeletePayload struct {
		ShelterId primitive.ObjectID   `json:"ShelterId" validate:"required"`
		PetsId    []primitive.ObjectID `json:"PetId" validate:"required"`
		UserId    primitive.ObjectID   `json:"UserId,omitempty"`
	}
	// Change Ready For Adopt Payload
	UpdateReadyForAdoptPayload struct {
		PetId     primitive.ObjectID `json:"PetId" validate:"required"`
		ShelterId primitive.ObjectID `json:"ShelterId" validate:"required"`
		UserId    primitive.ObjectID `json:"UserId" validate:"required"`
	}
)
