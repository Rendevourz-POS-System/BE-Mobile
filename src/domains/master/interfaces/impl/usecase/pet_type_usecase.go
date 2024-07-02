package usecase

import (
	"context"
	PetType "main.go/src/domains/master/entities"
	"main.go/src/domains/master/interfaces"
)

type petTypeUsecase struct {
	petRepo interfaces.PetTypeRepository
}

func NewPetTypeUsecase(pet interfaces.PetTypeRepository) *petTypeUsecase {
	return &petTypeUsecase{pet}
}

func (uc *petTypeUsecase) GetAllPetTypes(ctx context.Context) ([]PetType.PetType, error) {
	data, err := uc.petRepo.FindAllPets(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (uc *petTypeUsecase) CreatePetType(ctx context.Context, req *PetType.PetType) (*PetType.PetType, error) {
	data, err := uc.petRepo.StorePetType(ctx, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}
