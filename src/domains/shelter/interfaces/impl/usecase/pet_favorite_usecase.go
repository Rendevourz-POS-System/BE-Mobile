package usecase

import (
	"context"
	Shelter "main.go/domains/shelter/entities"
	"main.go/domains/shelter/interfaces"
)

type petFavoriteUseCase struct {
	petFavoriteRepo interfaces.PetFavoriteRepository
}

func NewPetFavoriteUseCase(petFavoriteRepo interfaces.PetFavoriteRepository) *petFavoriteUseCase {
	return &petFavoriteUseCase{petFavoriteRepo}
}

func (u *petFavoriteUseCase) UpdateIsFavoritePet(ctx context.Context, req *Shelter.PetFavoriteCreate) error {
	request := &Shelter.PetFavorite{
		PetId:  req.PetId,
		UserId: req.UserId,
	}
	if err := u.petFavoriteRepo.StoreOrUpdateIsFavoritePet(ctx, request); err != nil {
		return err
	}
	return nil
}
