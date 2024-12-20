package usecase

import (
	"context"
	Shelter "main.go/src/domains/shelter/entities"
	"main.go/src/domains/shelter/interfaces"
)

type shelterFavoriteUseCase struct {
	shelterFavoriteRepo interfaces.ShelterFavoriteRepository
}

func NewShelterFavoriteUseCase(shelterFavoriteRepo interfaces.ShelterFavoriteRepository) *shelterFavoriteUseCase {
	return &shelterFavoriteUseCase{shelterFavoriteRepo}
}

func (u *shelterFavoriteUseCase) UpdateIsFavoriteShelter(ctx context.Context, req *Shelter.ShelterFavoriteCreate) error {
	request := &Shelter.ShelterFavorite{
		ShelterId: req.ShelterId,
		UserId:    req.UserId,
	}
	if err := u.shelterFavoriteRepo.StoreOrUpdateIsFavorite(ctx, request); err != nil {
		return err
	}
	return nil
}
