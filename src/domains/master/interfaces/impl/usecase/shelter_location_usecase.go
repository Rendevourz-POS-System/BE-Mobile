package usecase

import (
	"context"
	ShelterLocation "main.go/domains/master/entities"
	"main.go/domains/master/interfaces"
)

type shelterLocationUsecase struct {
	shelterLocationRepo interfaces.ShelterLocationRepository
}

func NewShelterLocationUsecase(shelterLocationRepo interfaces.ShelterLocationRepository) *shelterLocationUsecase {
	return &shelterLocationUsecase{shelterLocationRepo}
}

func (u *shelterLocationUsecase) GetAllShelterLocation(ctx context.Context) ([]ShelterLocation.ShelterLocation, error) {
	data, err := u.shelterLocationRepo.FindAllShelterLocation(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *shelterLocationUsecase) CreateShelterLocation(ctx context.Context, req []ShelterLocation.ShelterLocation) (res []ShelterLocation.ShelterLocation, errs []string) {
	interfacesData := []interface{}{}
	for _, location := range req {
		interfacesData = append(interfacesData, location)
	}
	data, err := u.shelterLocationRepo.StoreShelterLocation(ctx, interfacesData)
	if err != nil {
		errs = append(errs, err.Error())
		return nil, errs
	}
	return data, nil
}
