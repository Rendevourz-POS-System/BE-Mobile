package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

type donationShelterRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewDonationShelterRepository(database *mongo.Database) *donationShelterRepo {
	return &donationShelterRepo{database: database, collection: database.Collection(collections.DonationShelterName)}
}
