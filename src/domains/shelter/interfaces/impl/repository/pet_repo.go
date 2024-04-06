package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

type petRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewPetRepository(database *mongo.Database) *petRepo {
	return &petRepo{database, database.Collection(collections.PetCollectionName)}
}
