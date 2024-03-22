package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	User      *mongo.Collection
	UserToken *mongo.Collection
)

func Migrate(db *mongo.Client, dbName string) error {
	User = db.Database(dbName).Collection("users")
	UserToken = db.Database(dbName).Collection("user_tokens")
	return nil
}
