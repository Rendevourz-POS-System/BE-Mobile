package usecase

import (
	"context"
	Shelter "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
	"main.go/domains/shelter/presistence"
	"main.go/shared/helpers"
)

type shelterUsecase struct {
	shelterRepo interfaces.ShelterRepository
}

func NewShelterUsecase(shelterRepo interfaces.ShelterRepository) *shelterUsecase {
	return &shelterUsecase{shelterRepo}
}

func (u *shelterUsecase) GetAllData(ctx context.Context, search *Shelter.ShelterSearch) ([]Shelter.Shelter, error) {
	if search.Sort == "" {
		search.Sort = "Desc"
	}
	if !presistence.ShelterFilterTable[search.OrderBy] || search.OrderBy == "" {
		search.OrderBy = "created_at"
	}
	data, err := u.shelterRepo.FindAllData(ctx, search)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *shelterUsecase) RegisterShelter(ctx context.Context, shelter *Shelter.Shelter) (*Shelter.Shelter, error) {
	shelter.CreatedAt = helpers.GetCurrentTime(nil)
	data, err := u.shelterRepo.StoreData(ctx, shelter)
	if err != nil {
		if data != nil {
			return data, err
		}
		return nil, err
	}
	return data, nil
}

func (u *shelterUsecase) GetOneDataById(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error) {
	data, err := u.shelterRepo.FindOneDataById(ctx, &search.ShelterId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *shelterUsecase) GetOneDataByUserId(ctx context.Context, search *Shelter.ShelterSearch) (*Shelter.Shelter, error) {
	data, err := u.shelterRepo.FindOneDataByUserId(ctx, &search.ShelterId)
	if err != nil {
		return nil, err
	}
	return data, nil
}