package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

type userTokenRepo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUserTokenRepository(database *mongo.Database) *userTokenRepo {
	return &userTokenRepo{database, database.Collection(collections.UserTokenCollectionName)}
}

func (repo *userTokenRepo) StoreToken(token string) error {
	return nil
}
