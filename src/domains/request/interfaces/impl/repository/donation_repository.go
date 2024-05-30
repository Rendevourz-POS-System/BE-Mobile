package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	Request "main.go/domains/request/entities"
	"main.go/shared/collections"
)

type donationShelterRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewDonationShelterRepository(database *mongo.Database) *donationShelterRepo {
	return &donationShelterRepo{database: database, collection: database.Collection(collections.DonationShelterName)}
}

func (repo *donationShelterRepo) StoreOneDonation(ctx context.Context, req *Request.DonationShelter) (res *Request.DonationShelter, err error) {
	data, errs := repo.collection.InsertOne(ctx, req)
	if errs != nil {
		return nil, errs
	}
	if err = repo.collection.FindOne(ctx, data).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
