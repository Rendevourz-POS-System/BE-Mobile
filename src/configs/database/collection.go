package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	UserPresistence "main.go/domains/user/presistence"
)

var (
	User      *mongo.Collection
	UserToken *mongo.Collection
)

func Migrate(db *mongo.Client, dbName string) error {
	User = db.Database(dbName).Collection(UserPresistence.UserCollectionName)
	UserToken = db.Database(dbName).Collection(UserPresistence.UserTokenCollectionName)
	return nil
}
