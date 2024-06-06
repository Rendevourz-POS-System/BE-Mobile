package interfaces

import (
	"context"
	Shelter "main.go/domains/shelter/entities"
)

type PetFavoriteUseCase interface {
	UpdateIsFavoritePet(ctx context.Context, req *Shelter.PetFavoriteCreate) error
}

type PetFavoriteRepository interface {
	StoreOrUpdateIsFavoritePet(ctx context.Context, req *Shelter.PetFavorite) error
}
