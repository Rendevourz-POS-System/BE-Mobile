package interfaces

import (
	"context"
	ShelterLocation "main.go/domains/master/entities"
)

type ShelterLocationRepository interface {
	FindAllShelterLocation(ctx context.Context) ([]ShelterLocation.ShelterLocation, error)
	StoreShelterLocation(ctx context.Context, req []interface{}) ([]ShelterLocation.ShelterLocation, error)
}

type ShelterLocationUsecase interface {
	GetAllShelterLocation(ctx context.Context) ([]ShelterLocation.ShelterLocation, error)
	CreateShelterLocation(ctx context.Context, req []ShelterLocation.ShelterLocation) ([]ShelterLocation.ShelterLocation, []string)
}
