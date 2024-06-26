package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	Pet "main.go/domains/shelter/entities"
	User "main.go/domains/user/entities"
	"mime/multipart"
	"time"
)

type Request struct {
	Id        primitive.ObjectID  `json:"Id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID  `json:"UserId" bson:"user_id"`
	PetId     *primitive.ObjectID `json:"PetId" bson:"pet_id"`
	ShelterId primitive.ObjectID  `json:"ShelterId" bson:"shelter_id"`
	Type      string              `json:"Type" bson:"type" validate:"required,request-type"`
	Status    string              `json:"Status" bson:"status" validate:"omitempty" default:"New"`
	//Job         string             `json:"Job" bson:"job"`
	Reason      *string    `json:"Reason,omitempty" bson:"reason"`
	RequestedAt *time.Time `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
	UpdatedAt   *time.Time `json:"UpdatedAt,omitempty" bson:"UpdatedAt,omitempty"`
	CompletedAt *time.Time `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
}

type (
	DonationPayload struct {
		Id        primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
		UserId    primitive.ObjectID `json:"UserId" bson:"user_id"`
		ShelterId primitive.ObjectID `json:"ShelterId" bson:"shelter_id" validate:"required"`
		RequestId primitive.ObjectID `json:"RequestId,omitempty" bson:"request_id"`
		Type      string             `json:"Type" bson:"type" validate:"required,donations"`
		Status    string             `json:"Status" bson:"status" validate:"omitempty" default:"New"`
		//Job         string             `json:"Job" bson:"job"`
		Reason         *string    `json:"Reason,omitempty" bson:"reason"`
		Amount         int64      `json:"Amount" validate:"omitempty"`
		PaymentType    string     `json:"PaymentType" bson:"payment_type" validate:"required,payment_type"`
		PaymentChannel *string    `json:"PaymentChannel,omitempty" bson:"payment_channel" validate:"required"`
		RequestedAt    *time.Time `json:"RequestedAt,omitempty" bson:"RequestedAt,omitempty"`
		CompletedAt    *time.Time `json:"CompletedAt,omitempty" bson:"CompletedAt,omitempty"`
	}
	RequestResponse struct {
		Request         *Request         `json:"Request,omitempty"`
		Donation        *DonationShelter `json:"Donation,omitempty"`
		Adoption        *AdoptionShelter `json:"Adoption,omitempty"`
		DonationPayload *DonationPayload `json:"DonationPayload,omitempty"`
		User            *User.User       `json:"User,omitempty"`
		UserTarget      *User.User       `json:"UserTarget,omitempty"`
	}
	SearchRequestPayload struct {
		RequestId *primitive.ObjectID `json:"RequestId,omitempty"`
		UserId    *primitive.ObjectID `json:"UserId,omitempty"`
		ShelterId *primitive.ObjectID `json:"ShelterId,omitempty"`
		Type      *[]string           `json:"Type"`
		Reason    *string             `json:"Reason"`
		Search    *string             `json:"Search"`
		Status    *string             `json:"Status"`
		Page      int                 `json:"Page"`
		PageSize  int                 `json:"PageSize"`
	}
	CreateRescueAndSurrenderRequestPayload struct {
		Files   *multipart.FileHeader `form:"Files" bson:"-" validate:"omitempty"`
		Pet     *Pet.Pet              `form:"Pet" bson:"Pet" validate:"required"`
		Request *Request              `form:"Request" bson:"-" validate:"required"`
	}
	RescueAndSurrenderResponse struct {
		Pet     *Pet.Pet `form:"Pet" bson:"Pet" validate:"required"`
		Request *Request `form:"Request" bson:"-" validate:"required"`
	}
	UpdateRescueAndSurrenderRequestStatus struct {
		RequestId primitive.ObjectID `json:"RequestId" bson:"request_id" validate:"required"`
		Type      string             `json:"Type" bson:"type" validate:"required,request-type,rescueOrSurrender"`
		Status    string             `json:"Status" bson:"status" validate:"required,approve_or_reject"`
		Reason    *string            `json:"Reason,omitempty" bson:"reason"`
	}
	UpdateRescueAndSurrenderRequestStatusResponse struct {
		Request *Request `json:"Request"`
		Pet     *Pet.Pet `json:"Pet"`
	}
	UpdateAdoptionRequestStatus struct {
		RequestId primitive.ObjectID `json:"RequestId" validate:"required"`
		ShelterId primitive.ObjectID `json:"ShelterId"`
	}
)
