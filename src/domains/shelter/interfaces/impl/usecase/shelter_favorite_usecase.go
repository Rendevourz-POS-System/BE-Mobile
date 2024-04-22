package usecase

import (
	"context"
	Shelter "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
)

type shelterFavoriteUseCase struct {
	shelterFavoriteRepo interfaces.ShelterFavoriteRepository
}

func NewShelterFavoriteUseCase(shelterFavoriteRepo interfaces.ShelterFavoriteRepository) *shelterFavoriteUseCase {
	return &shelterFavoriteUseCase{shelterFavoriteRepo}
}

func (u *shelterFavoriteUseCase) UpdateIsFavoriteShelter(ctx context.Context, req *Shelter.ShelterFavoriteCreate) error {
	if err := u.shelterFavoriteRepo.StoreOrUpdateIsFavorite(ctx, req); err != nil {
		return err
	}
	return nil
}
