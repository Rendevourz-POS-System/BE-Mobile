package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	ShelterLocation "main.go/domains/master/entities"
	"main.go/shared/collections"
)

type shelterLocationRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewShelterLocationRepository(database *mongo.Database) *shelterLocationRepo {
	return &shelterLocationRepo{database, database.Collection(collections.ShelterLocationName)}
}

func (r *shelterLocationRepo) FindAllShelterLocation(ctx context.Context) (res []ShelterLocation.ShelterLocation, err error) {
	data, err := r.collection.Find(ctx, bson.M{})
	if err = data.All(ctx, &res); err != nil {

	}
	return res, nil
}
