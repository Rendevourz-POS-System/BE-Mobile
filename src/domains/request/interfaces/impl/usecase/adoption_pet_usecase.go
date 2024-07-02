package usecase

import "main.go/src/domains/request/interfaces"

type adoptionPetUsecase struct {
	adoptionPetRepo interfaces.AdoptionPetRepository
}

func NewAdoptionPetUsecase(adoptionPet interfaces.AdoptionPetRepository) *adoptionPetUsecase {
	return &adoptionPetUsecase{adoptionPet}
}
