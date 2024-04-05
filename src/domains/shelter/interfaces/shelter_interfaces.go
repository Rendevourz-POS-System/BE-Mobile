package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Shelter "main.go/domains/shelter/entities"
)

type ShelterRepository interface {
	FindAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.Shelter, error)
	StoreData(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, error)
	FindOneDataById(c context.Context, search *primitive.ObjectID) (res *Shelter.Shelter, err error)
	FindOneDataByUserId(c context.Context, search *primitive.ObjectID) (res *Shelter.Shelter, err error)
}

type ShelterUsecase interface {
	GetAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.Shelter, error)
	RegisterShelter(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, error)
	GetOneDataById(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error)
	GetOneDataByUserId(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error)
}
