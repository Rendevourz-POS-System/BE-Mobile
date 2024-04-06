package interfaces

import Pet "main.go/domains/shelter/entities"

type PetRepository interface {
}

type PetUseCase interface {
	GetAllPets(search *Pet.PetSearch) ([]Pet.Pet, error)
}
