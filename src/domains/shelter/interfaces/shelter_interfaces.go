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
	UpdatePet(ctx context.Context, pet *Shelter.Shelter) (*Shelter.Shelter, error)
}

type ShelterUsecase interface {
	GetAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.Shelter, error)
	RegisterShelter(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, []string)
	GetOneDataById(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error)
	GetOneDataByUserId(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error)
	UpdatePetById(ctx context.Context, Id *primitive.ObjectID, search *Shelter.Shelter) (*Shelter.Shelter, error)
}
