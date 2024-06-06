package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	Shelter "main.go/domains/shelter/entities"
	"main.go/shared/collections"
)

type petFavoriteRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewPetFavoriteRepository(database *mongo.Database) *petFavoriteRepo {
	return &petFavoriteRepo{database, database.Collection(collections.PetFavoriteName)}
}

func (r *petFavoriteRepo) StoreOrUpdateIsFavoritePet(ctx context.Context, req *Shelter.PetFavorite) error {
	// Try to find a document matching the given criteria
	foundDoc := r.collection.FindOne(ctx, bson.M{"pet_id": req.PetId, "user_id": req.UserId})
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
		_, err := r.collection.DeleteOne(ctx, bson.M{"pet_id": req.PetId, "user_id": req.UserId})
		if err != nil {
			return err
		}
	}
	return nil
}
