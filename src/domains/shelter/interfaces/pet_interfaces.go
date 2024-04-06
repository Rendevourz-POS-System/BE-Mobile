package interfaces

import (
	"context"
	Pet "main.go/domains/shelter/entities"
)

type PetRepository interface {
	FindAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.Pet, error)
}

type PetUseCase interface {
	GetAllPets(ctx context.Context, search *Pet.PetSearch) ([]Pet.Pet, error)
}
