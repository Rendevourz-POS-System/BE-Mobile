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
	City               string             `json:"City" bson:"city" validate:"required"`
	PostalCode         int                `json:"PostalCode" bson:"postal_code" validate:"required,number"`
	Province           string             `json:"Province" bson:"province" validate:"required"`
	Email              string             `json:"Email" bson:"email" validate:"required,email"`
	Username           string             `json:"Username" bson:"username" validate:"required,min=4"`
	Password           string             `json:"Password" bson:"password" validate:"required,min=8,alphanum_symbol"`
	StaffStatus        bool               `json:"StaffStatus" bson:"staff_status" default:"false" validate:"omitempty"`
	ShelterIsActivated bool               `json:"ShelterIsActivated" bson:"shelter_is_activated" default:"false" validate:"omitempty"`
	Role               string             `json:"Role" bson:"role" validate:"omitempty,required,role" default:"User"`
	Verified           bool               `json:"Verified" bson:"is_active"`
	CreatedAt          *time.Time         `json:"createdAt" bson:"CreatedAt,omitempty"`
	DeletedAt          *time.Time         `json:"deletedAt,omitempty" bson:"DeletedAt,omitempty"`
}

type (
	// LoginPayload Payload for login
	LoginPayload struct {
		Email    string `json:"Email" validate:"required,email"`
		Username string `json:"Username" validate:"omitempty"`
		Password string `json:"Password" validate:"required"`
	}
	// LoginResponse Response for login
	LoginResponse struct {
		Username string `json:"Username"`
		Token    string `json:"Token"`
	}
	// JwtCustomClaims Custom claims for JWT
	JwtCustomClaims struct {
		ID       string `json:"Id"`
		Email    string `json:"Email"`
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
		Nonce string `json:"Nonce"`
		jwt.RegisteredClaims
	}
)
