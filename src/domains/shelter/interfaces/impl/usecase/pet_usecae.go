package usecase

import (
	"context"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
)

type petUseCase struct {
	petRepo interfaces.PetRepository
}

func NewPetUseCase(petRepo interfaces.PetRepository) *petUseCase {
	return &petUseCase{petRepo}
}

func (u *petUseCase) GetAllPets(ctx context.Context, search *Pet.PetSearch) (res []Pet.Pet, err error) {
	if res, err = u.petRepo.FindAllPets(ctx, search); err != nil {
		return nil, err
	}
	return res, nil
}
