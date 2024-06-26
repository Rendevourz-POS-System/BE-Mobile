package usecase

import (
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/configs/app"
	Pet "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/shared/helpers"
	"main.go/shared/helpers/image_helpers"
	"os"
)

type petUseCase struct {
	petRepo interfaces.PetRepository
}

func NewPetUseCase(petRepo interfaces.PetRepository) *petUseCase {
	return &petUseCase{petRepo}
}

func (u *petUseCase) CheckIsValidForUpdate(ctx context.Context, Id *primitive.ObjectID) (bool, error) {
	return u.petRepo.ValidateIfValidForUpdate(ctx, Id)
}

func (u *petUseCase) GetAllPets(ctx context.Context, search *Pet.PetSearch) (res []Pet.PetResponsePayload, err error) {
	if res, err = u.petRepo.FindAllPets(ctx, search); err != nil {
		return nil, err
	}
	for i, pet := range res {
		var base64Images []string
		for _, imagePath := range pet.Image {
			var path string
			if pet.ShelterId != nil && len(pet.ShelterId.Hex()) > 0 {
				path = image_helpers.GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, pet.ShelterId.Hex(), app.GetConfig().Image.PetPath, pet.ID.Hex(), imagePath)
			} else {
				path = image_helpers.GenerateImagePath(app.GetConfig().Image.PetPath, pet.ID.Hex(), imagePath)
			}
			imageData, err := os.ReadFile(path) // Read the image file
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
	flag := false
	if pet.Pet.IsAdopted == nil {
		pet.Pet.IsAdopted = &flag
	}
	if pet.Pet.ReadyToAdopt == nil {
		pet.Pet.ReadyToAdopt = &flag
	}
	pet.Pet.CreatedAt = helpers.GetCurrentTime(nil)
	if res, err = u.petRepo.StorePets(ctx, &pet.Pet); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *petUseCase) UpdatePet(ctx context.Context, Id *primitive.ObjectID, pet *Pet.Pet) (res *Pet.Pet, err error) {
	pet.ID = *Id
	flag := false
	if pet.ReadyToAdopt == nil {
		pet.ReadyToAdopt = &flag
	}
	if pet.IsAdopted == nil {
		pet.IsAdopted = &flag
	}
	data, errs := u.petRepo.UpdatePet(ctx, pet)
	if errs != nil {
		return nil, errs
	}
	if len(data.Image) > 0 {
		for _, imagePath := range data.Image {
			var path string
			if data.ShelterId != nil && len(data.ShelterId.Hex()) > 0 {
				path = image_helpers.GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, pet.ShelterId.Hex(), app.GetConfig().Image.PetPath, pet.ID.Hex(), imagePath)
			} else {
				path = image_helpers.GenerateImagePath(app.GetConfig().Image.PetPath, pet.ID.Hex(), imagePath)
			}
			imageData, err := os.ReadFile(path) // Read the image file
			if err != nil {
				return nil, err // Handle error (perhaps just log and continue with other images?)
			}
			base64Image := base64.StdEncoding.EncodeToString(imageData) // Convert to Base64
			data.ImageBase64 = append(data.ImageBase64, base64Image)
		}
	}
	return data, nil
}

func (u *petUseCase) GetPetById(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error) {
	res, err = u.petRepo.FindPetById(ctx, Id)
	if err != nil {
		return nil, err
	}
	for _, imagePath := range res.Image {
		var path string
		if res.ShelterId != nil && len(res.ShelterId.Hex()) > 0 {
			path = image_helpers.GenerateImagePath(app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, res.ShelterId.Hex(), app.GetConfig().Image.PetPath, res.ID.Hex(), imagePath)
		} else {
			path = image_helpers.GenerateImagePath(app.GetConfig().Image.PetPath, res.ID.Hex(), imagePath)
		}
		imageData, err := os.ReadFile(path) // Read the image file
		if err != nil {
			return nil, err // Handle error (perhaps just log and continue with other images?)
		}
		base64Image := base64.StdEncoding.EncodeToString(imageData) // Convert to Base64
		res.ImageBase64 = append(res.ImageBase64, base64Image)
	}
	return res, nil
}

func (u *petUseCase) DeletePetByAdmin(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error) {
	res, err = u.petRepo.DestroyPetByAdmin(ctx, Id)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *petUseCase) DeletePetByUser(ctx context.Context, pet Pet.PetDeletePayload) (res []Pet.Pet, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(pet); errs != nil {
		err = helpers.CustomError(errs)
		return nil, err
	}
	res, err = u.petRepo.DestroyPetByUser(ctx, pet)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *petUseCase) UpdateReadyForAdoptStatus(ctx context.Context, Id *primitive.ObjectID) (res *Pet.Pet, err error) {
	res, err = u.petRepo.PutReadyForAdoptStatus(ctx, Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
