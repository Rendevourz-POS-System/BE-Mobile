package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/src/shared/collections"
)

var (
	User             *mongo.Collection
	UserToken        *mongo.Collection
	PetType          *mongo.Collection
	ShelterFavorites *mongo.Collection
	Shelter          *mongo.Collection
	ShelterLocation  *mongo.Collection
	Request          *mongo.Collection
	Adoption         *mongo.Collection
	Donation         *mongo.Collection
)

func Migrate(db *mongo.Client, dbName string) error {
	User = db.Database(dbName).Collection(collections.UserCollectionName)
	UserToken = db.Database(dbName).Collection(collections.UserTokenCollectionName)
	Shelter = db.Database(dbName).Collection(collections.ShelterCollectionName)
	PetType = db.Database(dbName).Collection(collections.PetTypeName)
	ShelterFavorites = db.Database(dbName).Collection(collections.ShelterFavoriteName)
	ShelterLocation = db.Database(dbName).Collection(collections.ShelterLocationName)
	Request = db.Database(dbName).Collection(collections.RequestName)
	Adoption = db.Database(dbName).Collection(collections.Adoption)
	Donation = db.Database(dbName).Collection(collections.DonationShelterName)
	return nil
}
