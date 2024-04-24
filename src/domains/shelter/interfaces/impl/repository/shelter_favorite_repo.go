package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *shelterFavoriteRepo) StoreOrUpdateIsFavorite(ctx context.Context, req *Shelter.ShelterFavorite) error {
	// Try to find a document matching the given criteria
	foundDoc := r.collection.FindOne(ctx, bson.M{"shelter_id": req.ShelterId, "user_id": req.UserId})
	if foundDoc.Err() != nil {
		if foundDoc.Err() == mongo.ErrNoDocuments {
			// No document found, insert a new one
			_, err := r.collection.InsertOne(ctx, req)
			if err != nil {
				return err
			}
		} else {
			return foundDoc.Err()
		}
	} else {
		// Document found, delete it
		_, err := r.collection.DeleteOne(ctx, bson.M{"shelter_id": req.ShelterId, "user_id": req.UserId})
		if err != nil {
			return err
		}
	}
	return nil
}
