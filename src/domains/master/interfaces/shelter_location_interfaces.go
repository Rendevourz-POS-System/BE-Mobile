package interfaces

import (
	"context"
	ShelterLocation "main.go/domains/master/entities"
)

type ShelterLocationRepository interface {
	FindAllShelterLocation(ctx context.Context) ([]ShelterLocation.ShelterLocation, error)
}

type ShelterLocationUsecase interface {
	GetAllShelterLocation(ctx context.Context) ([]ShelterLocation.ShelterLocation, error)
}
