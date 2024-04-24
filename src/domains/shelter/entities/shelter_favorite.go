package entities

type ShelterFavorite struct {
	ShelterId string `json:"ShelterId" bson:"shelter_id" validate:"required"`
	UserId    string `json:"UserId" bson:"user_id" validate:"required"`
}

type (
	// ShelterFavoriteCreate Payload
	ShelterFavoriteCreate struct {
		UserId    string `json:"UserId"`
		ShelterId string `json:"ShelterId" validate:"required"`
	}
)
