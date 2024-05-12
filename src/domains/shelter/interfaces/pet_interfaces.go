package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Pet "main.go/domains/shelter/entities"
)

type PetRepository interface {
	FindAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.PetResponsePayload, error)
	StorePets(ctx context.Context, pet *Pet.Pet) (*Pet.Pet, []string)
	UpdatePet(ctx context.Context, pet *Pet.Pet) (*Pet.Pet, error)
}

type PetUseCase interface {
	GetAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.PetResponsePayload, error)
	CreatePets(ctx context.Context, pet *Pet.PetCreate) (*Pet.Pet, []string)
	UpdatePet(ctx context.Context, Id *primitive.ObjectID, pet *Pet.Pet) (res *Pet.Pet, err error)
}
