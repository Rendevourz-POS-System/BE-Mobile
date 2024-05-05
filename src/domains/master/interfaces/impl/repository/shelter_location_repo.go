package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return nil, err
	}
	return res, nil
}

func (r *shelterLocationRepo) StoreShelterLocation(ctx context.Context, req []interface{}) (res []ShelterLocation.ShelterLocation, err error) {
	data, err := r.collection.InsertMany(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Handle the duplicate key error
			return nil, fmt.Errorf("duplicate key error: %v", err)
		}
		return nil, err
	}
	// Prepare a slice to collect _id's of the inserted documents
	var ids []primitive.ObjectID
	for _, id := range data.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			ids = append(ids, oid)
		}
	}
	// Query to find all inserted documents
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
