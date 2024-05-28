package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AdoptionShelter struct {
	Id                  primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	RequestId           primitive.ObjectID `json:"RequestId" bson:"request_id" validate:"required"`
	PetId               primitive.ObjectID `json:"PetId" bson:"pet_id" validate:"required"`
	DateApproval        *time.Time         `json:"DateApproval,omitempty" bson:"date_approval"`
	MonitoringPhase     int8               `json:"MonitoringPhase,omitempty" bson:"monitoring_phase"`
	MonitoringPetStatus *string            `json:"MonitoringPetStatus,omitempty" bson:"monitoring_pet_status"`
}
