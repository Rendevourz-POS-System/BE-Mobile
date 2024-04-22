package usecase

import (
	"context"
	"fmt"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/shared/helpers"
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

func (u *petUseCase) CreatePets(ctx context.Context, pet *Pet.PetCreate) (res []Pet.Pet, err []string) {
	validate := helpers.NewValidator()
	// Validate data
	if errs := validate.Struct(pet); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	if pet.Pet.PetGender == "" {
		pet.Pet.PetGender = "Unknown"
	}
	pet.Pet.CreatedAt = helpers.GetCurrentTime(nil)
	if res, err = u.petRepo.StorePets(ctx, &pet.Pet); err != nil {
		return nil, err
	}
	fmt.Println("petData", pet.Pet)
	return res, nil
}
