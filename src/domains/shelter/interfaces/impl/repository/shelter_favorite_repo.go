package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	Shelter "main.go/domains/shelter/entities"
	"main.go/shared/collections"
)

type shelterFavoriteRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewShelterFavoriteRepository(database *mongo.Database) *shelterFavoriteRepo {
	return &shelterFavoriteRepo{database, database.Collection(collections.ShelterFavoriteName)}
}

func (r *shelterFavoriteRepo) StoreOrUpdateIsFavorite(ctx context.Context, req *Shelter.ShelterFavoriteCreate) error {
	_, err := r.collection.InsertOne(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
