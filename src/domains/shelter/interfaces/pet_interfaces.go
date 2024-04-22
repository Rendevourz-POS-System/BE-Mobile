package interfaces

import (
	"context"
	Pet "main.go/domains/shelter/entities"
)

type PetRepository interface {
	FindAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.Pet, error)
	StorePets(ctx context.Context, pet *Pet.Pet) ([]Pet.Pet, []string)
}

type PetUseCase interface {
	GetAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.Pet, error)
	CreatePets(ctx context.Context, pet *Pet.PetCreate) ([]Pet.Pet, []string)
}
