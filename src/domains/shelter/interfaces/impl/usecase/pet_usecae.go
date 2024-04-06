package usecase

import (
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
)

type petUseCase struct {
	petRepo interfaces.PetRepository
}

func NewPetUseCase(petRepo interfaces.PetRepository) *petUseCase {
	return &petUseCase{petRepo}
}

func (u *petUseCase) GetAllPets(search *Pet.PetSearch) ([]Pet.PetSearch, error) {

	return nil, nil
}
