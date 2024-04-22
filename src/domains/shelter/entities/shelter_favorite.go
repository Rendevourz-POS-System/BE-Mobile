package entities

type ShelterFavorite struct {
	ShelterId string `json:"shelter_id" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
}

type (
	// ShelterFavoriteCreate Payload
	ShelterFavoriteCreate struct {
		ShelterId string `json:"shelter_id" validate:"required"`
	}
)
