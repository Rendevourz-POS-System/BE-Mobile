package usecase

import (
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/src/configs/app"
	Shelter "main.go/src/domains/shelter/entities"
	"main.go/src/domains/shelter/interfaces"
	"main.go/src/domains/shelter/presistence"
	"main.go/src/shared/helpers"
	"main.go/src/shared/helpers/image_helpers"
	"os"
	"time"
)

type shelterUsecase struct {
	shelterRepo interfaces.ShelterRepository
}

func NewShelterUsecase(shelterRepo interfaces.ShelterRepository) *shelterUsecase {
	return &shelterUsecase{shelterRepo}
}

func (u *shelterUsecase) GetAllData(ctx context.Context, search *Shelter.ShelterSearch) (res []Shelter.ShelterResponsePayload, err error) {
	if search.Sort == "" {
		search.Sort = "Desc"
	}
	if !presistence.ShelterFilterTable[search.OrderBy] || search.OrderBy == "" {
		search.OrderBy = "created_at"
	}
	res, err = u.shelterRepo.FindAllData(ctx, search)
	if err != nil {
		return nil, err
	}
	for i, item := range res {
		var base64Images []string
		for _, imagePath := range item.Image {
			imageData, err := os.ReadFile(image_helpers.GenerateImagePath(
				app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, item.ID.Hex(), imagePath)) // Read the image file
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

func (u *shelterUsecase) RegisterShelter(ctx context.Context, shelter *Shelter.Shelter) (res *Shelter.Shelter, err []string) {
	validate := helpers.NewValidator()
	if errs := validate.Struct(shelter); errs != nil {
		err := helpers.CustomError(errs)
		return nil, err
	}
	shelter.CreatedAt = helpers.GetCurrentTime(nil)
	data, errs := u.shelterRepo.StoreData(ctx, shelter)
	if errs != nil {
		err = append(err, errs.Error())
		if data != nil {
			return data, err
		}
		return nil, err
	}
	return data, nil
}

func (u *shelterUsecase) GetOneDataById(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.ShelterResponsePayload, error) {
	data, err := u.shelterRepo.FindOneDataById(ctx, &search.ShelterId)
	if err != nil {
		return nil, err
	}
	var base64Images []string
	for _, imagePath := range data.Image {
		imageData, err := os.ReadFile(image_helpers.GenerateImagePath(
			app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, data.ID.Hex(), imagePath)) // Read the image file
		if err != nil {
			return nil, err // Handle error (perhaps just log and continue with other images?)
		}
		base64Image := base64.StdEncoding.EncodeToString(imageData) // Convert to Base64
		base64Images = append(base64Images, base64Image)
	}
	data.ImageBase64 = base64Images // Assuming pets have an ImageBase64 field to store the base64 strings
	return data, nil
}

func (u *shelterUsecase) GetOneDataByIdForRequest(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error) {
	data, err := u.shelterRepo.FindOneDataByIdForRequest(ctx, &search.ShelterId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *shelterUsecase) GetOneDataByUserId(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.ShelterResponsePayload, error) {
	data, err := u.shelterRepo.FindOneDataByUserId(ctx, &search.UserId)
	if err != nil {
		return nil, err
	}
	for _, imagePath := range data.Image {
		imageData, err := os.ReadFile(image_helpers.GenerateImagePath(
			app.GetConfig().Image.UserPath, app.GetConfig().Image.ShelterPath, data.ID.Hex(), imagePath)) // Read the image file
		if err != nil {
			return nil, err // Handle error (perhaps just log and continue with other images?)
		}
		base64Image := base64.StdEncoding.EncodeToString(imageData) // Convert to Base64
		data.ImageBase64 = append(data.ImageBase64, base64Image)
	}
	return data, nil
}

func (u *shelterUsecase) UpdateShelterById(ctx context.Context, Id *primitive.ObjectID, shelter *Shelter.Shelter) (res *Shelter.Shelter, err error) {
	shelter.ID = *Id
	data, errs := u.shelterRepo.UpdateOneShelter(ctx, shelter)
	if errs != nil {
		return nil, errs
	}
	return data, nil
}

func (u *shelterUsecase) DeleteAllDataShelterByAdmin(ctx context.Context, shelterId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := u.shelterRepo.DestroyAllDataShelterByAdmin(ctx, shelterId); err != nil {
		return err
	}
	return nil
}
