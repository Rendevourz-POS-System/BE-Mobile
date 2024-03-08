package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	User *mongo.Collection
)

func Migrate(db *mongo.Client, dbName string) error {
	User = db.Database(dbName).Collection("users")
	return nil
}
