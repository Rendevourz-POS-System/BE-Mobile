package usecase

import (
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/shared/helpers"
	"os"
)

type petUseCase struct {
	petRepo interfaces.PetRepository
}

func NewPetUseCase(petRepo interfaces.PetRepository) *petUseCase {
	return &petUseCase{petRepo}
}

func (u *petUseCase) GetAllPets(ctx context.Context, search *Pet.PetSearch) (res []Pet.PetResponsePayload, err error) {
	if res, err = u.petRepo.FindAllPets(ctx, search); err != nil {
		return nil, err
	}
	for i, pet := range res {
		var base64Images []string
		for _, imagePath := range pet.Image {
			imageData, err := os.ReadFile(imagePath) // Read the image file
			if err != nil {
				return nil, err // Handle error (perhaps just log and continue with other images?)
			}
			base64Image := base64.StdEncoding.EncodeToString(imageData) // Convert to Base64
			base64Images = append(base64Images, base64Image)
		}
		res[i].ImageBase64 = base64Images // Assuming pets have an ImageBase64 field to store the base64 strings
	}
	return res, nil
}

func (u *petUseCase) CreatePets(ctx context.Context, pet *Pet.PetCreate) (res *Pet.Pet, err []string) {
	validate := helpers.NewValidator()
	// Validate data
	if errs := validate.Struct(pet); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	if pet.Pet.PetGender == "" {
		pet.Pet.PetGender = "Unknown"
	}
	pet.Pet.PetDob = nil
	pet.Pet.CreatedAt = helpers.GetCurrentTime(nil)
	if res, err = u.petRepo.StorePets(ctx, &pet.Pet); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *petUseCase) UpdatePet(ctx context.Context, Id *primitive.ObjectID, pet *Pet.Pet) (res *Pet.Pet, err error) {
	pet.ID = *Id
	data, errs := u.petRepo.UpdatePet(ctx, pet)
	if errs != nil {
		return nil, errs
	}
	return data, nil
}

func (u *petUseCase) GetPetById(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error) {
	res, err = u.petRepo.FindPetById(ctx, Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
