package interfaces

import (
	"context"
	PetType "main.go/src/domains/master/entities"
)

type PetTypeRepository interface {
	FindAllPets(ctx context.Context) ([]PetType.PetType, error)
	StorePetType(ctx context.Context, req *PetType.PetType) (*PetType.PetType, error)
}

type PetTypeUsecase interface {
	GetAllPetTypes(ctx context.Context) ([]PetType.PetType, error)
	CreatePetType(ctx context.Context, req *PetType.PetType) (*PetType.PetType, error)
}
