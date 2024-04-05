package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/shared/collections"
)

var (
	User      *mongo.Collection
	UserToken *mongo.Collection
)

func Migrate(db *mongo.Client, dbName string) error {
	User = db.Database(dbName).Collection(collections.UserCollectionName)
	UserToken = db.Database(dbName).Collection(collections.UserTokenCollectionName)
	return nil
}
