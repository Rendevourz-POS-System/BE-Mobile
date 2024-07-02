package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Shelter "main.go/src/domains/shelter/entities"
)

type ShelterRepository interface {
	FindAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.ShelterResponsePayload, error)
	StoreData(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, error)
	FindOneDataById(c context.Context, search *primitive.ObjectID) (res *Shelter.ShelterResponsePayload, err error)
	FindOneDataByUserId(c context.Context, search *primitive.ObjectID) (res *Shelter.ShelterResponsePayload, err error)
	UpdateOneShelter(ctx context.Context, pet *Shelter.Shelter) (*Shelter.Shelter, error)
	FindOneDataByIdForRequest(c context.Context, Id *primitive.ObjectID) (res *Shelter.Shelter, err error)
}

type ShelterUsecase interface {
	GetAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.ShelterResponsePayload, error)
	RegisterShelter(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, []string)
	GetOneDataById(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.ShelterResponsePayload, error)
	GetOneDataByIdForRequest(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error)
	GetOneDataByUserId(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.ShelterResponsePayload, error)
	UpdateShelterById(ctx context.Context, Id *primitive.ObjectID, search *Shelter.Shelter) (*Shelter.Shelter, error)
}
