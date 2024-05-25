package entities

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Entities
type User struct {
	ID                 primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	Nik                string             `json:"Nik" bson:"nik" validate:"required"`
	PhoneNumber        string             `json:"PhoneNumber" bson:"phone_number" validate:"required,number"`
	Address            string             `json:"Address" bson:"address" validate:"required"`
	State              string             `json:"State" bson:"state" default:"Indonesia"`
	City               string             `json:"City" bson:"city" validate:"required"`
	Province           string             `json:"Province" bson:"province" validate:"required"`
	District           string             `json:"District" bson:"district" validate:"required"`
	PostalCode         int                `json:"PostalCode" bson:"postal_code" validate:"required,number"`
	Email              string             `json:"Email" bson:"email" validate:"required,email"`
	Username           string             `json:"Username" bson:"username" validate:"required,min=4"`
	Password           string             `json:"Password" bson:"password" validate:"required,min=8,alphanum_symbol"`
	StaffStatus        bool               `json:"StaffStatus" bson:"staff_status" default:"false" validate:"omitempty"`
	ShelterIsActivated bool               `json:"ShelterIsActivated" bson:"shelter_is_activated" default:"false" validate:"omitempty"`
	Role               string             `json:"Role" bson:"role" validate:"omitempty,required,role" default:"User"`
	Image              string             `json:"Image" bson:"image" validate:"omitempty"`
	ImageBase64        string             `json:"ImageBase64" validate:"omitempty"`
	Verified           bool               `json:"Verified" bson:"is_active"`
	CreatedAt          *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
	UpdatedAt          *time.Time         `json:"UpdatedAt" bson:"UpdatedAt,omitempty"`
	DeletedAt          *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	// Update Profile Payload
	UpdateProfilePayload struct {
		ID                 primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
		Nik                string             `json:"Nik" bson:"nik" validate:"required"`
		PhoneNumber        string             `json:"PhoneNumber" bson:"phone_number" validate:"required,number"`
		Address            string             `json:"Address" bson:"address" validate:"required"`
		State              string             `json:"State" bson:"state" validate:"required"`
		City               string             `json:"City" bson:"city" validate:"required"`
		Province           string             `json:"Province" bson:"province" validate:"required"`
		District           string             `json:"District" bson:"district" validate:"required"`
		PostalCode         int                `json:"PostalCode" bson:"postal_code" validate:"required,number"`
		Email              string             `json:"Email" bson:"email" validate:"required,email"`
		Username           string             `json:"Username" bson:"username" validate:"required,min=4"`
		StaffStatus        bool               `json:"StaffStatus" bson:"staff_status" default:"false" validate:"omitempty"`
		ShelterIsActivated bool               `json:"ShelterIsActivated" bson:"shelter_is_activated" validate:"omitempty"`
		Role               string             `json:"Role" bson:"role" validate:"omitempty,required,role"`
		Image              string             `json:"Image,omitempty" bson:"image" validate:"omitempty"`
		OldImageName       string             `json:"OldImageName,omitempty" validate:"omitempty"`
		ImageBase64        string             `json:"ImageBase64" validate:"omitempty"`
		Verified           bool               `json:"Verified" bson:"is_active"`
		CreatedAt          *time.Time         `json:"CreatedAt" bson:"CreatedAt,omitempty"`
		UpdatedAt          *time.Time         `json:"UpdatedAt" bson:"UpdatedAt,omitempty"`
		DeletedAt          *time.Time         `json:"DeletedAt,omitempty" bson:"DeletedAt,omitempty"`
	}
	// Update Password Payload
	UpdatePasswordPayload struct {
		Id              primitive.ObjectID `json:"Id"`
		Password        string             `json:"Password" bson:"password" validate:"required,min=8,alphanum_symbol"`
		NewPassword     string             `json:"NewPassword" validate:"required,min=8,alphanum_symbol"`
		ConfirmPassword string             `json:"ConfirmPassword" validate:"required,eqfield=NewPassword"`
	}
	// Verfied Email Payload
	EmailVerifiedPayload struct {
		Token  string             `json:"Token"`
		UserId primitive.ObjectID `json:"UserId" validate:"required"`
		Otp    *int               `json:"Otp" validate:"required"`
	}
	// LoginPayload Payload for login
	LoginPayload struct {
		Email    string `json:"Email" validate:"required,email"`
		Username string `json:"Username" validate:"omitempty"`
		Password string `json:"Password" validate:"required"`
	}
	// LoginResponse Response for login
	LoginResponse struct {
		User     User   `json:"User,omitempty"`
		Username string `json:"Username"`
		Token    string `json:"Token"`
	}
	// JwtCustomClaims Custom claims for JWT
	JwtCustomClaims struct {
		ID       string `json:"Id"`
		Email    string `json:"Email"`
		Otp      *int   `json:"Otp"`
		Username string `json:"Username"`
		jwt.RegisteredClaims
	}
	// JwtCustomRefreshClaims Custom claims for JWT Refresh Token
	JwtCustomRefreshClaims struct {
		ID string `json:"Id"`
		jwt.RegisteredClaims
	}
	// JwtEmailClaims Custom claims for JWT Email Verification
	JwtEmailClaims struct {
		ID    string `json:"Id"`
		Email string `json:"Email"`
		Otp   *int   `json:"Otp"`
		Nonce string `json:"Nonce"`
		jwt.RegisteredClaims
	}
)
