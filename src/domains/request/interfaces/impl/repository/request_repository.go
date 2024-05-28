package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

type requestRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewRequestRepository(database *mongo.Database) *requestRepo {
	return &requestRepo{database: database, collection: database.Collection(collections.RequestName)}
}
