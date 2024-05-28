package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

type adoptionPetRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewAdoptionPetRepository(database *mongo.Database) *adoptionPetRepo {
	return &adoptionPetRepo{database: database, collection: database.Collection(collections.AdoptionPetName)}
}
