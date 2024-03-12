package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	User "main.go/domains/user/entities"
	"main.go/domains/user/presistence"
)

type userRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *userRepository {
	return &userRepository{database, database.Collection(presistence.CollectionName)}
}

func (userRepo *userRepository) FindAll(c context.Context) (res []User.User, err error) {
	data, err := userRepo.collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	err = data.All(c, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (userRepo *userRepository) StoreOne(c context.Context, user *User.User) (res *User.User, errs error) {
	insertResult, err := userRepo.collection.InsertOne(c, user)
	if err != nil {
		return nil, err
	}

	// Ensure the InsertedID is an ObjectID to use it in the FindOne query
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err // You might want to return a more descriptive error here
	}

	if err = userRepo.collection.FindOne(c, bson.M{"_id": objectID}).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
