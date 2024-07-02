package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Pet "main.go/src/domains/shelter/entities"
)

type PetRepository interface {
	FindAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.PetResponsePayload, error)
	StorePets(ctx context.Context, pet *Pet.Pet) (*Pet.Pet, []string)
	UpdatePet(ctx context.Context, pet *Pet.Pet) (*Pet.Pet, error)
	PutReadyForAdoptStatus(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error)
	FindPetById(ctx context.Context, Id *primitive.ObjectID) (*Pet.Pet, error)
	DestroyPetByAdmin(ctx context.Context, Id *primitive.ObjectID) (*Pet.Pet, error)
	DestroyPetByUser(ctx context.Context, Pets Pet.PetDeletePayload) ([]Pet.Pet, []string)
	ValidateIfValidForUpdate(ctx context.Context, Id *primitive.ObjectID) (bool, error)
}

type PetUseCase interface {
	GetAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.PetResponsePayload, error)
	CreatePets(ctx context.Context, pet *Pet.PetCreate) (*Pet.Pet, []string)
	UpdatePet(ctx context.Context, Id *primitive.ObjectID, pet *Pet.Pet) (res *Pet.Pet, err error)
	UpdateReadyForAdoptStatus(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error)
	GetPetById(ctx context.Context, Id *primitive.ObjectID) (*Pet.Pet, error)
	DeletePetByAdmin(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error)
	DeletePetByUser(ctx context.Context, Pets Pet.PetDeletePayload) (res []Pet.Pet, err []string)
	CheckIsValidForUpdate(ctx context.Context, Id *primitive.ObjectID) (bool, error)
}
