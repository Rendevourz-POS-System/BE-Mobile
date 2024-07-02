package interfaces

import (
	"context"
	Shelter "main.go/src/domains/shelter/entities"
)

type ShelterFavoriteUseCase interface {
	UpdateIsFavoriteShelter(ctx context.Context, req *Shelter.ShelterFavoriteCreate) error
}

type ShelterFavoriteRepository interface {
	StoreOrUpdateIsFavorite(ctx context.Context, req *Shelter.ShelterFavorite) error
}
