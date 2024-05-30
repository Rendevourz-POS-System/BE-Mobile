package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	Request "main.go/domains/request/entities"
	"main.go/shared/collections"
)

type requestRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewRequestRepository(database *mongo.Database) *requestRepo {
	return &requestRepo{database: database, collection: database.Collection(collections.RequestName)}
}

func (r *requestRepo) StoreOneRequest(ctx context.Context, req *Request.Request) (res *Request.Request, err error) {
	data, errInsert := r.collection.InsertOne(ctx, req)
	if errInsert != nil {
		return nil, errInsert
	}
	if err = r.collection.FindOne(ctx, bson.M{"_id": data.InsertedID}).Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
